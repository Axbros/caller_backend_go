package handler

import (
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
			dataValue, exists := dataMap["data"]
			if exists {
				eventStr, ok := eventValue.(string)
				if ok {
					err := w.iDao.SetIPAddrByMachineCode2WebsocketConnections(ctx, dataValue.(string), remoteAddr)
					if err != nil {
						logger.Errorf("set connection error %s", eventStr, logger.Err(err))
						return
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

		logger.Infof("get websocket client message:%s from %s", messageStr, remoteAddr)
	}
}
