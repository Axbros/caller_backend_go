package handler

import (
	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/model"
	"caller/internal/utils"
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/ws"
)

var clients = make(map[string]*websocket.Conn)

type WebsocketHandler interface {
	LoopReceiveMessage(ctx context.Context, conn *ws.Conn)
}
type websocketHandler struct {
	iDao dao.RedisDao
}

func NewWebsocketHandler() WebsocketHandler {
	return &websocketHandler{
		iDao: dao.NewRedisDao(model.GetRedisCli()),
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
					dataStr, _ := dataMap["data"].(string)
					userTypeStr, _ := dataMap["type"].(string)
					if eventStr == "heartbeat" {
						clients[dataStr] = conn
						updateHeartBeatInfo(w, ctx, dataStr, remoteAddr)
					}
					if userTypeStr == "client" { //处理客户端的消息
						if eventStr == "income" || eventStr == "endcall" || eventStr == "connected" {
							userDao := dao.NewUserDao(
								model.GetDB(),
								cache.NewUserCache(model.GetCacheType()),
							)
							parent, err := userDao.GetUserMachineCodeByClientMachineCode(ctx, dataStr)
							parentMachineCode := ""
							if err != nil {
								logger.Errorf("GetUserMachineCodeByClientMachineCode error", err)
							}
							//无论话机是否绑定甲方都将数据存储在redis种
							dataMap["from"] = remoteAddr
							jsonStr, err := json.Marshal(dataMap)

							if len(parent) > 0 {
								parentMachineCode = parent[0].MachineCode
								dataMap["to"] = parentMachineCode
								//把指令放在redis存储 当receive方收到之后执行Delete操作
								redisStoreKey := parentMachineCode + ":" + dataMap["key"].(string)
								w.iDao.SetMessageStore(ctx, redisStoreKey, jsonStr)
								if err != nil {
									logger.Error("存储指令到redis种失败", logger.Err(err))
									return
								}
								sendDataToSpecificClient(clients[parentMachineCode], jsonStr)
								logger.Info("websocket", logger.String("send message", messageStr))
							} else {
								err := conn.WriteMessage(websocket.TextMessage, []byte("当前没有在线的甲方设备"))
								logger.Info("Not Found", logger.String("client", remoteAddr), logger.String("user", parentMachineCode))
								if err != nil {
									logger.Error("send message error", logger.Err(err))
								}
							}

							// {"event":"income","message":"16600229957","data":"689550a77428cee9","from":"192.168.1.220:38002"}
							// {"event":"endcall","message":"16600229957","data":"689550a77428cee9"} from 192.168.1.220:38002
						}
					}

					if userTypeStr == "user" {
						//处理用户端的消息
						if eventStr == "receive" { //表示用户端收到话机的指令 需要执行清除redis操作
							redisStoreKey := dataStr + ":" + dataMap["message"].(string) // 88888888:testkey
							w.iDao.DeleteMessageStore(ctx, redisStoreKey)
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

func sendDataToSpecificClient(conn *ws.Conn, message []byte) {

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		logger.Error("向客户端发送数据出错", logger.Err(err), logger.String("message", string(message)), logger.String("to", conn.RemoteAddr().String()))
		err := conn.Close()
		if err != nil {
			return
		}
	}
}
