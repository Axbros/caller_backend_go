package handler

import (
	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/model"
	"caller/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/ws"
)

var clients = make(map[string]*websocket.Conn)

type WebsocketHandler interface {
	LoopReceiveMessage(ctx context.Context, conn *ws.Conn)
}
type websocketHandler struct {
	iDao dao.RedisDao
	dDao dao.DistributionDao
	gDao dao.GroupCallDao
	uDao dao.UserDao
}

func NewWebsocketHandler() WebsocketHandler {
	return &websocketHandler{
		iDao: dao.NewRedisDao(model.GetRedisCli()),
		dDao: dao.NewDistributionDao(model.GetDB(),
			cache.NewDistributionCache(model.GetCacheType())),
		gDao: dao.NewGroupCallDao(model.GetDB(),
			cache.NewGroupCallCache(model.GetCacheType())),
		uDao: dao.NewUserDao(model.GetDB(),
			cache.NewUserCache(model.GetCacheType())),
	}
}

func (w websocketHandler) LoopReceiveMessage(ctx context.Context, conn *ws.Conn) {

	for {
		_, message, err := conn.ReadMessage()
		remoteAddr := conn.RemoteAddr().String()
		if err != nil {
			log.Println("ReadMessage error:", err)
			break
		}
		// 将字节切片转换为字符串
		messageStr := string(message)
		// 处理消息
		jsonData, err := utils.ParseTextToJSON(messageStr)

		if err != nil {
			logger.Error("ParseTextToJSON error", logger.Err(err), logger.String("origin message is", messageStr))
			return
		}
		// 类型断言判断是否为 map[string]interface{} 类型
		dataMap, ok := jsonData.(map[string]interface{})
		if ok {
			eventValue, exists := dataMap["event"]

			if exists {
				eventStr, ok := eventValue.(string)
				if ok {
					jsonStr, _ := json.Marshal(dataMap)
					dataStr, _ := dataMap["data"].(string)
					userTypeStr, _ := dataMap["type"].(string)
					if eventStr == "heartbeat" {
						clients[dataStr] = conn
						updateHeartBeatInfo(w, ctx, dataStr, remoteAddr)
					}

					if userTypeStr == "client" { //处理客户端的消息
						if eventStr == "income" || eventStr == "endcall" || eventStr == "connected" || eventStr == "call_done" {
							parent, err := w.uDao.GetUserByClientMachineCode(ctx, dataStr)
							if err != nil {
								logger.Error("GetUserByClientMachineCode error", logger.Err(err))
							}
							//无论话机是否绑定甲方都将数据存储在redis种
							dataMap["from"] = remoteAddr
							parentIdStr := strconv.FormatUint(parent.ID, 10)
							if parent.ID > 0 {
								dataMap["to"] = parentIdStr
								//把指令放在redis存储 当receive方收到之后执行Delete操作
								redisStoreKey := parentIdStr + ":" + dataMap["key"].(string)
								w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
								if err != nil {
									logger.Error("存储指令到redis种失败", logger.Err(err))
									return
								}
								parentConn := clients[parentIdStr]
								sendDataToSpecificClient(parentConn, jsonStr)
								logger.Info("websocket", logger.String("send message", messageStr), logger.String("receive address", parentConn.RemoteAddr().String()))
							} else {
								err := conn.WriteMessage(websocket.TextMessage, []byte("当前没有在线的甲方设备,甲方ID："+parentIdStr))
								logger.Info("Not Found", logger.String("client", remoteAddr), logger.String("user", parentIdStr))
								if err != nil {
									logger.Error("send message error", logger.Err(err))
								}
							}
						}

						if eventStr == "transfer_done" {
							client_id := dataMap["data"].(string)
							messageKey := dataMap["key"].(string)
							//此时中专机应在data将自己的id传过来 然后用这个id去group_call表查询transfer_client_id 为 {id}的记录
							groupcall_record, err := w.gDao.GetByCondition(ctx, &query.Conditions{
								Columns: []query.Column{
									{
										Name:  "transfer_client_id",
										Value: client_id,
									},
								},
							})
							if err != nil {
								logger.Err(err)
							}
							//查询到上方的记录就有了group_call.id 拿到这个id去distribution表查询group_call_id=group_call.id的记录
							distribution_record, _ := w.dDao.GetByCondition(ctx, &query.Conditions{
								Columns: []query.Column{
									{
										Name:  "group_call_id",
										Value: groupcall_record.ID,
									},
								}})

							//上方查询成功后会拿到user_id 把配置成功的消息发给user

							sendDataToSpecificClient(clients[strconv.Itoa(distribution_record.UserID)], generateServerWebsocketMsg("中转机已成功配置，即将拨打电话", messageKey))
							//开始给话机拨打电话
							queen_client, _ := w.iDao.GetQueenValue(ctx, "group_name_"+groupcall_record.GroupName)
							sendDataToSpecificClient(clients[strconv.Itoa(distribution_record.UserID)], generateServerWebsocketMsg("您的话机组名为：【"+groupcall_record.GroupName+"】已从队列取出话机："+queen_client+"即将呼出目标号码", messageKey))
							//给话机传递指令 拨打中转号码
							sendDataToSpecificClient(clients[queen_client], generateStandardWebsocketMsg("call", groupcall_record.PhoneNumber, messageKey))
						}
					}

					if userTypeStr == "user" {
						//处理用户端的消息
						logger.Info("receive from user", logger.String("message", messageStr), logger.String("event", eventStr))
						data := dataMap["data"].(string)
						messageStr := dataMap["message"].(string)
						messageKey := dataMap["key"].(string)
						redisStoreKey := data + ":" + messageKey // 88888888:testkey
						if eventStr == "receive" {               //表示用户端收到话机的指令 需要执行清除redis操作
							w.iDao.DeleteMessageStore(ctx, redisStoreKey)
						}
						if eventStr == "missed" || eventStr == "outgoing" {
							children := strings.Split(messageStr, ",")

							for _, child := range children {
								if clients[child] != nil {
									sendDataToSpecificClient(clients[child], jsonStr)
								}
							}
						}
						if eventStr == "endcall" {
							w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
							//user excute endcall or answer should put the client machine code id to message,data is user machine code id
							sendDataToSpecificClient(clients[messageStr], jsonStr)
							//这里的data其实就是client machine code
							//话机执行了挂断操作 需要 把结果告诉甲方
							sendDataToSpecificClient(conn, generateServerWebsocketMsg("decline", messageKey))
						}
						if eventStr == "answer" {
							// w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
							// //user excute endcall or answer should put the client machine code id to message,data is user machine code id
							sendDataToSpecificClient(clients[messageStr], jsonStr)
							return
						}
						if eventStr == "call" {
							//这里面的data就是本机userID
							w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
							sendDataToSpecificClient(conn, generateServerWebsocketMsg("服务端收到指令，正在配置中转设备", messageKey))
							//todo 根据userID查询中转设备
							group := w.dDao.GetDistributedGroupCallIdByUserId(ctx, data)
							if group > 0 {
								sendDataToSpecificClient(conn, generateServerWebsocketMsg("查询到已采用中转方案，正在获取中转信息", messageKey))
								group_call_record, err := w.gDao.GetByID(ctx, group)
								if err != nil {
									logger.Error("get transfer record error", logger.Err(err), logger.String("group id", string(group)))
								}
								transfer_phone := group_call_record.PhoneNumber
								transfer_machine_id := group_call_record.TransferClientID
								sendDataToSpecificClient(conn, generateServerWebsocketMsg(fmt.Sprintf("中转号码：%s 中转设备:%s", transfer_phone, transfer_machine_id), messageKey))
								if clients[transfer_machine_id] != nil {
									err := sendDataToSpecificClient(clients[transfer_machine_id], generateStandardWebsocketMsg("transfer", messageStr, messageKey))
									if err != nil {
										sendDataToSpecificClient(conn, generateServerWebsocketMsg("中转配置出错，中转设备ID："+transfer_machine_id, messageKey))
									}
								} else {
									// for key, value := range clients {
									// 	fmt.Printf("Key: %s, Value: %s\n", key, value.RemoteAddr)
									// }
									sendDataToSpecificClient(conn, generateServerWebsocketMsg("中转设备不在线！请联系管理员处理，中转设备ID："+transfer_machine_id, messageKey))
								}

							} else {
								sendDataToSpecificClient(conn, generateServerWebsocketMsg("当前没有分配中转设备，即将选择直射模式", messageKey))
							}
						}
					}

				} else {
					logger.Error("event 的值不是字符串类型")
				}
			} else {
				logger.Error("event 键不存在")
			}
		} else {
			logger.Error("解析结果不是预期的 map 类型")
		}
		logger.Info("websocket", logger.String("messageStr", messageStr), logger.String("remoteAddr", remoteAddr))
	}
}

func updateHeartBeatInfo(w websocketHandler, ctx context.Context, machine_code, ip_address string) {
	err := w.iDao.SetIPAddrByMachineCode2WebsocketConnections(ctx, machine_code, ip_address)
	if err != nil {
		logger.Errorf("set connection error %s", "heartbeat", logger.Err(err))
		return
	}
}

func sendDataToSpecificClient(conn *ws.Conn, message []byte) error {
	logger.Info("websocket", logger.String("send to", conn.RemoteAddr().String()), logger.String("message", string(message)))
	if conn != nil {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.Error("向客户端发送数据出错", logger.Err(err), logger.String("message", string(message)), logger.String("to", conn.RemoteAddr().String()))
			err := conn.Close()
			if err != nil {
				return err
			}
		}
	} else {
		logger.Error("接收设备不在线:" + conn.RemoteAddr().String())
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func generateServerWebsocketMsg(message, key string) []byte {
	msg := fmt.Sprintf(`{"event":"receive","message":"%s","data":"","key":"%s","type":"server"}`, message, key)
	return []byte(msg)
}

func generateStandardWebsocketMsg(event, message, key string) []byte {
	msg := fmt.Sprintf(`{"event":"%s","message":"%s","data":"","key":"%s","type":"server"}`, event, message, key)
	return []byte(msg)
}
