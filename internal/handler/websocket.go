package handler

import (
	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/model"
	"caller/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/ws"
)

var HEARTBEAT_TIME float64 = 65
var rwMu sync.RWMutex
var clients = make(map[string]*websocket.Conn)
var ip2deviceID = make(map[string]string)

func deleteClient(key string) {
	rwMu.Lock()
	delete(clients, key)
	rwMu.Unlock() // 添加Unlock来释放锁
	logger.Info("删除设备", logger.Any("设备ID", key))
}

func readFromClients(key string) *websocket.Conn {
	rwMu.RLock()
	defer rwMu.RUnlock() // 使用defer确保锁被释放
	value := clients[key]
	return value
}

func updateClients(key string, value *websocket.Conn) {
	logger.Info("检测到设备加入", logger.Any("设备ID", key))
	ip2deviceID[value.RemoteAddr().String()] = key
	rwMu.Lock()
	clients[key] = value
	rwMu.Unlock() // 确保在记录日志前释放锁
}

type WebsocketHandler interface {
	LoopReceiveMessage(ctx context.Context, conn *ws.Conn)
	GetOnlineClients(ctx *gin.Context)
}
type websocketHandler struct {
	iDao dao.RedisDao
	dDao dao.DistributionDao
	gDao dao.GroupCallDao
	uDao dao.UserDao
	cDao dao.UnanswerdCallDao
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
		cDao: dao.NewUnanswerdCallDao(model.GetDB(),
			cache.NewUnanswerdCallCache(model.GetCacheType())),
	}
}
func closeConn(remoteAddr string) {
	offlineDeviceId := ip2deviceID[remoteAddr]
	delete(ip2deviceID, remoteAddr)
	deleteClient(offlineDeviceId)
}
func appendOfflinePhoneNumber(ctx context.Context, parentId string, phoneNumber string) {
	redisDao := dao.NewRedisDao(model.GetRedisCli())
	// 把phoneNumber 添加到redis队列中
	redisDao.PushOfflinePhoneNumber(context.Background(), "offline:"+parentId, phoneNumber)
}
func getOfflinePhoneNumber(ctx context.Context, parentId string) ([]string, error) {
	redisDao := dao.NewRedisDao(model.GetRedisCli())
	return redisDao.GetAllOfflinePhoneNumber(ctx, "offline:"+parentId)
}
func (w websocketHandler) LoopReceiveMessage(ctx context.Context, conn *ws.Conn) {

	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	lastHeartbeatTime := time.Time{} // 初始化上一次收到心跳的时间为零值
	// 注册关闭事件处理函数
	conn.SetCloseHandler(func(code int, text string) error {
		closeConn(remoteAddr)
		logger.Info("WebSocket客户端断开连接", logger.Any("code", code), logger.Any("reason", text), logger.Any("还剩设备:", len(clients)))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Info("读取WebSocket消息出错", logger.Any("err", err), logger.Any("出错地址", remoteAddr))
			return
		}

		// 将字节切片转换为字符串
		messageStr := string(message)
		// 处理消息
		jsonData, err := utils.ParseTextToJSON(messageStr)

		if err != nil {
			logger.Error("ParseTextToJSON error", logger.Err(err), logger.String("origin message is", messageStr))
			continue
		}
		// 类型断言判断是否为 map[string]interface{} 类型
		dataMap, ok := jsonData.(map[string]interface{})
		if ok {
			eventValue, exists := dataMap["event"]

			if exists {
				eventStr, ok := eventValue.(string)
				if ok {
					jsonStr, _ := json.Marshal(dataMap)
					dataStr, _ := dataMap["data"].(string) // 话机ID
					userTypeStr, _ := dataMap["type"].(string)
					if eventStr == "heartbeat" {
						currentTime := time.Now()
						if lastHeartbeatTime.IsZero() {
							// 如果是首次收到心跳消息，记录当前时间
							lastHeartbeatTime = currentTime
							if userTypeStr == "user" { //甲方上线 首先检查是否有离线消息
								//read from redis offline queue
								offlinePhoneNumbers, err := getOfflinePhoneNumber(ctx, dataStr)
								if err != nil {
									logger.Error("GetAllOfflinePhoneNumber error", logger.Err(err))
								}
								sendDataToSpecificClient(conn, generateStandardWebsocketMsg("offline", strings.Join(offlinePhoneNumbers, ","), "", "offline"))
							}
						} else {
							// todo 话机没有发送心跳 配置好了心跳再打开下面代码
							// 计算距离首次收到心跳的时间间隔
							duration := currentTime.Sub(lastHeartbeatTime)
							if duration.Seconds() > HEARTBEAT_TIME {
								logger.Info("等待设备发送心跳超时，主动断开")
								closeConn(remoteAddr)
								// 时间间隔大于10秒，断开连接
								conn.Close()
								break
							}
						}
						// clients[dataStr] = conn
						// 更新上一次收到心跳的时间为当前时间
						lastHeartbeatTime = currentTime
						logger.Info("收到心跳包", logger.Any("设备ID", dataStr), logger.Any("地址", remoteAddr), logger.Any("当前时间", currentTime))
						updateClients(dataStr, conn)
						// reply client that the sever has recieve the message
						// sendDataToSpecificClient(conn, generateServerWebsocketMsg("I am sever,I have recived your message", "hi"))
						// updateHeartBeatInfo(w, ctx, dataStr, remoteAddr)

					}

					if userTypeStr == "client" { //处理客户端的消息

						if eventStr == "income" || eventStr == "endcall" || eventStr == "connected" || eventStr == "call_done" {

							parent, err := w.uDao.GetUserByClientMachineCode(ctx, dataStr)
							if err != nil {
								logger.Error("GetUserByClientMachineCode error", logger.Err(err))
							}
							//无论话机是否绑定甲方都将数据存储在redis种
							// dataMap["from"] = remoteAddr
							parentIdStr := strconv.FormatUint(parent.ID, 10)
							if parent.ID > 0 {
								dataMap["to"] = parentIdStr
								parentConn := readFromClients(parentIdStr)
								if parentConn != nil {
									sendDataToSpecificClient(parentConn, jsonStr)
								} else {
									appendOfflinePhoneNumber(ctx, parentIdStr, dataMap["message"].(string))
								}
								if eventStr == "income" {
									w.cDao.Create(ctx, &model.UnanswerdCall{
										MachineId:    dataStr,
										MobileNumber: dataMap["message"].(string),
										Location:     dataMap["from"].(string),
										Type:         "income",
									})
								}

							} else {
								// 甲方设备不在线 存在缓存中 等待上线后进行推送

								logger.Info("未绑定甲方设备", logger.String("client", remoteAddr), logger.String("user", parentIdStr))
							}
						}

						if eventStr == "transfer_done" {
							client_id := dataMap["data"].(string)
							messageKey := dataMap["key"].(string)
							message := dataMap["message"].(string)
							to := dataMap["to"].(string)
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
							userConn := readFromClients(strconv.Itoa(distribution_record.UserID))
							sendDataToSpecificClient(userConn, generateServerWebsocketMsg("中转机已成功配置，即将拨打电话", messageKey))
							//开始给话机拨打电话
							var targetMachine *websocket.Conn

							if to != "" {
								//指定话机拨打
								targetMachine = readFromClients(to)

							} else {
								for i := 0; i < len(clients); i++ {
									queen_client, _ := w.iDao.GetQueenValue(ctx, "group_name_"+groupcall_record.GroupName)
									if readFromClients(queen_client) != nil {
										sendDataToSpecificClient(readFromClients(strconv.Itoa(distribution_record.UserID)), generateServerWebsocketMsg("您的话机组名为：【"+groupcall_record.GroupName+"】已从队列取出话机："+queen_client+"即将呼出目标号码", messageKey))
										targetMachine = readFromClients(queen_client)
										w.cDao.Create(ctx, &model.UnanswerdCall{
											MachineId:    queen_client,
											MobileNumber: message,
											Location:     "中国·大陆",
											Type:         "keypad",
										},
										)
										break
									} else {
										sendDataToSpecificClient(readFromClients(strconv.Itoa(distribution_record.UserID)), generateServerWebsocketMsg("队列话机不在线，话机ID："+queen_client+"即将队列循环到下一个话机", messageKey))
										continue
									}
								}
							}

							if targetMachine != nil {
								//给话机传递指令 拨打中转号码 只有拨号盘拨打出的电话才同步到数据库

								sendDataToSpecificClient(targetMachine, generateStandardWebsocketMsg("call", groupcall_record.PhoneNumber, "", messageKey))

							} else {
								sendDataToSpecificClient(readFromClients(strconv.Itoa(distribution_record.UserID)), generateServerWebsocketMsg("队列话机不在线", messageKey))
							}

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
						// if eventStr == "missed" || eventStr == "outgoing" {
						// 	children := strings.Split(messageStr, ",")

						// 	for _, child := range children {
						// 		if readFromClients(child) != nil {
						// 			sendDataToSpecificClient(readFromClients(child), jsonStr)
						// 		}
						// 	}
						// }
						if eventStr == "endcall" {
							// w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
							//user excute endcall or answer should put the client machine code id to message,data is user machine code id
							sendDataToSpecificClient(readFromClients(messageStr), jsonStr)
							//这里的data其实就是client machine code
							//话机执行了挂断操作 需要 把结果告诉甲方
							sendDataToSpecificClient(conn, generateServerWebsocketMsg("decline", messageKey))
						}
						if eventStr == "answer" {
							// w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
							// //user excute endcall or answer should put the client machine code id to message,data is user machine code id
							sendDataToSpecificClient(readFromClients(messageStr), jsonStr)
							// locationStr := dataMap["from"].(string)
							w.cDao.Create(ctx, &model.UnanswerdCall{
								MachineId:    dataStr,
								MobileNumber: dataMap["message"].(string),
								Location:     dataMap["from"].(string),
								Type:         "answer",
							})
						}
						if eventStr == "call" {
							to := dataMap["to"].(string)
							//这里面的data就是本机userID
							// w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
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
								if transfer_machine_id != "0" {
									sendDataToSpecificClient(conn, generateServerWebsocketMsg(fmt.Sprintf("中转号码：%s 中转设备:%s", transfer_phone, transfer_machine_id), messageKey))

									if readFromClients(transfer_machine_id) != nil {
										err := sendDataToSpecificClient(readFromClients(transfer_machine_id), generateStandardWebsocketMsg("transfer", messageStr, to, messageKey))
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
									group_name := group_call_record.GroupName
									for i := 0; i < len(clients); i++ {
										queen_client, _ := w.iDao.GetQueenValue(ctx, "group_name_"+group_name)
										if readFromClients(queen_client) != nil {
											w.cDao.Create(ctx, &model.UnanswerdCall{
												MachineId:    queen_client,
												MobileNumber: messageStr,
												Type:         "keypad",
											},
											)
											sendDataToSpecificClient(readFromClients(queen_client), generateStandardWebsocketMsg("call", messageStr, "", messageKey))

											sendDataToSpecificClient(conn, generateStandardWebsocketMsg("read_success", "设备读取成功", queen_client, messageKey))

											break
										} else {
											sendDataToSpecificClient(conn, generateServerWebsocketMsg("队列话机不在线，话机ID："+queen_client+"即将队列循环到下一个话机", messageKey))
											continue
										}
									}
								}
								sendDataToSpecificClient(conn, generateStandardWebsocketMsg("flow_done", "流程执行结束", "", messageKey))
							} else {
								sendDataToSpecificClient(conn, generateStandardWebsocketMsg("404", "当前未查询到您的绑定关系，请联系管理员核实。", "", "oh_no"))
								// sendDataToSpecificClient(conn, generateStandardWebsocketMsg("查询到已采用中转方案，正在获取中转信息", messageKey))
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

func (w websocketHandler) GetOnlineClients(c *gin.Context) {
	res := make(map[string]string)
	for key, value := range clients {
		ipAddr := value.RemoteAddr().String()
		res[key] = ipAddr
	}
	response.Success(c, gin.H{
		"results": res,
		"count":   len(res),
	})
}

// func updateHeartBeatInfo(w websocketHandler, ctx context.Context, machine_code, ip_address string) {
// 	err := w.iDao.SetIPAddrByMachineCode2WebsocketConnections(ctx, machine_code, ip_address)
// 	if err != nil {
// 		logger.Errorf("set connection error %s", "heartbeat", logger.Err(err))
// 		return
// 	}
// }

func sendDataToSpecificClient(conn *ws.Conn, message []byte) error {
	if conn != nil {
		logger.Infof("开始给客户端发送消息", logger.Any("设备ID", conn.RemoteAddr().String()))
		if err := writeMessageWithLogging(conn, websocket.TextMessage, message); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond) // 暂停一段时间，可能为了防止消息发送过快
		return nil
	}
	return nil
}

func generateServerWebsocketMsg(message, key string) []byte {
	msg := fmt.Sprintf(`{"event":"receive","message":"%s","data":"","key":"%s","type":"server"}`, message, key)
	return []byte(msg)
}

func generateStandardWebsocketMsg(event, message, data, key string) []byte {

	msg := fmt.Sprintf(`{"event":"%s","message":"%s","data":"%s","key":"%s","type":"server"}`, event, message, data, key)
	return []byte(msg)
}
func writeMessageWithLogging(conn *ws.Conn, messageType int, message []byte) error {
	if conn != nil {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			logger.Error("向客户端发送数据出错", logger.Err(err), logger.String("message", string(message)), logger.String("to", conn.RemoteAddr().String()))
			// 尝试关闭连接，忽略关闭连接时可能出现的错误
			_ = conn.Close()
			return err
		}
	} else {
		logger.Error("接收设备不在线")
		return errors.New("device received is offline, the message has been aborted")
	}
	return nil
}
