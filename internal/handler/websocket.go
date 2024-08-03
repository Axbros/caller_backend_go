package handler

import (
	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/model"
	"caller/internal/utils"
	"context"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/ws"
	"log"
)

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
			logger.Error("ParseTextToJSON error", logger.Err(err))
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
						updateHeartBeatInfo(w, ctx, dataStr, remoteAddr)
					}
					if eventStr == "income" || eventStr == "endcall" || eventStr == "connected" {
						userDao := dao.NewUserDao(
							model.GetDB(),
							cache.NewUserCache(model.GetCacheType()),
						)
						parentMachineCode, err := userDao.GetUserMachineCodeByClientMachineCode(ctx, dataStr)
						if err != nil {
							logger.Errorf("GetUserMachineCodeByClientMachineCode error", err)
						}
						logger.Info("websocket", logger.String("UserMachineCode", parentMachineCode[0].MachineCode), logger.String("ClientMachineCode", dataStr))
						// {"event":"income","message":"16600229957","data":"689550a77428cee9"} from 192.168.1.220:38002
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
