package handler

import (
	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/model"
	"caller/internal/utils"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/ws"
	"log"
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
					if eventStr == "heartbeat" {
						clients[dataStr] = conn
						//sendDataToSpecificClient([]byte("rotatorw"))
						updateHeartBeatInfo(w, ctx, dataStr, remoteAddr)
					}
					if eventStr == "income" || eventStr == "endcall" || eventStr == "connected" {
						userDao := dao.NewUserDao(
							model.GetDB(),
							cache.NewUserCache(model.GetCacheType()),
						)
						parent, err := userDao.GetUserMachineCodeByClientMachineCode(ctx, dataStr)
						if err != nil {
							logger.Errorf("GetUserMachineCodeByClientMachineCode error", err)
						}
						//无论话机是否绑定甲方都将数据存储在redis种
						jsonStr, err := json.Marshal(dataMap)
						err = w.iDao.SetMessageStore(ctx, dataMap["key"].(string), jsonStr)
						if len(parent) > 0 {
							parentMachineCode := parent[0].MachineCode
							dataMap["from"] = remoteAddr
							dataMap["to"] = parentMachineCode
							//把指令放在redis存储 当receive方收到之后执行Delete操作

							if err != nil {
								logger.Error("存储指令到redis种失败", logger.Err(err))
								return
							}
							sendDataToSpecificClient(clients[parentMachineCode], messageStr)
							logger.Info("websocket", logger.String("send message", messageStr))
						} else {
							err := conn.WriteMessage(websocket.TextMessage, []byte("当前没有在线的话机"))
							if err != nil {
								logger.Error("send message error", logger.Err(err))
							}
						}

						// {"event":"income","message":"16600229957","data":"689550a77428cee9","from":"192.168.1.220:38002"}
						// {"event":"endcall","message":"16600229957","data":"689550a77428cee9"} from 192.168.1.220:38002
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

func sendDataToSpecificClient(conn *ws.Conn, message string) {
	messageByte := []byte(message)
	err := conn.WriteMessage(websocket.TextMessage, messageByte)
	if err != nil {
		logger.Error("向客户端发送数据出错", logger.Err(err), logger.String("message", message), logger.String("to", conn.RemoteAddr().String()))
		err := conn.Close()
		if err != nil {
			return
		}
	}
}
