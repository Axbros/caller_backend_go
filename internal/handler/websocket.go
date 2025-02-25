package handler

import (
	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/ws"
)

type ContentType string

// 定义 ContentType 的枚举值
const (
	Call        ContentType = "call"
	EndCall     ContentType = "endcall"
	Idle        ContentType = "idle"
	Ping        ContentType = "ping"
	Pong        ContentType = "pong"
	SendFail    ContentType = "SendFail"
	SendSuccess ContentType = "SendSuccess"
	SendSMS     ContentType = "SendSMS"
	Income      ContentType = "income"
	Answer      ContentType = "answer"
	NewSMS      ContentType = "new_sms"
)

// 定义一个结构体来存储连接信息
type ConnectionInfo struct {
	Conn              *websocket.Conn
	RemoteAddr        string    // 远程地址
	LastHeartbeatTime time.Time // 最后心跳时间
}

// 定义一个全局的 map 来存储所有用户的连接信息
var connectionMap = make(map[string]ConnectionInfo)

// 使用互斥锁来保证并发安全
var mutex sync.Mutex

// ContentBody 定义与 JavaScript 对象对应的结构体
type ContentBody struct {
	ContentType ContentType `json:"contentType"`
	Content     interface{} `json:"content"`
	SenderID    string      `json:"senderID"`
	ReceiveID   string      `json:"receiveID"`
	ClientMsgID string      `json:"clientMsgID"`
}

var rwMu sync.RWMutex
var clients = make(map[string]*websocket.Conn)

// 存储每个客户端的最后心跳时间
var clientLastHeartbeat = make(map[string]time.Time)

const (
	// 心跳间隔
	HeartbeatInterval = 5 * time.Second
	// 超时时间
	HeartbeatTimeout = 10 * time.Second
)

type WebsocketHandler interface {
	LoopReceiveMessage(ctx context.Context, conn *ws.Conn)
	GetOnlineClients(ctx *gin.Context)
	CheckHeartBeat()
}

type websocketHandler struct {
	iDao dao.GroupClientDao
}

func NewWebsocketHandler() WebsocketHandler {
	return &websocketHandler{
		iDao: dao.NewGroupClientDao(model.GetDB(),
			cache.NewGroupClientCache(model.GetCacheType())),
	}
}

func startHeartbeatCheck(done chan struct{}) {
	logger.Info("开始心跳检测")
	go func() {
		ticker := time.NewTicker(HeartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rwMu.RLock()
				for deviceID, lastTime := range clientLastHeartbeat {
					logger.Info("正在遍历客户端心跳数据", logger.Any("deviceID", deviceID), logger.Any("lastTime", lastTime))
					if time.Since(lastTime) > HeartbeatTimeout {
						rwMu.RUnlock()
						rwMu.Lock()
						fmt.Printf("客户端 %s 可能断网\n", deviceID)
						delete(connectionMap, deviceID)
						delete(clientLastHeartbeat, deviceID)
						rwMu.Unlock()
						rwMu.RLock()
					}
				}
				rwMu.RUnlock()
			case <-done:
				return
			}
		}
	}()
}

func (w websocketHandler) CheckHeartBeat() {
	logger.Info("开始监听心跳包")
	done := make(chan struct{})
	startHeartbeatCheck(done)
}

// 处理消息发送和错误处理的通用函数
func handleMessageSend(conn *ws.Conn, messageType int, contentBody ContentBody, action ContentType) error {
	err := sendMessageToUser(conn, contentBody)
	if err != nil {
		message := createMessage(SendFail, action, contentBody.SenderID, contentBody.ClientMsgID)
		return conn.WriteMessage(messageType, message)
	}
	return nil
}

// 获取群主 ID 的通用函数
func getGroupOwnerID(ctx context.Context, w websocketHandler, senderID string) (string, error) {
	groupOwnerID, err := w.iDao.GetGroupOwnerIDByClientID(ctx, senderID)
	if err != nil {
		logger.Error("GetGroupNameByClientID error", logger.Err(err), logger.Any("clientID", senderID))
		return "", err
	}
	return groupOwnerID, nil
}

func (w websocketHandler) LoopReceiveMessage(ctx context.Context, conn *ws.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()

	conn.SetCloseHandler(func(code int, text string) error {
		logger.Info("WebSocket客户端断开连接", logger.Any("code", code), logger.Any("reason", text), logger.Any("还剩设备:", len(clientLastHeartbeat)))
		return nil
	})

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Info("读取WebSocket消息出错", logger.Any("err", err), logger.Any("出错地址", remoteAddr))
			// 根据remoteAddr删除clientLastHeartbeat
			var deviceID string
			for k, v := range connectionMap {
				if v.RemoteAddr == remoteAddr {
					deviceID = k
					break
				}
			}
			logger.Info("删除设备", logger.Any("设备ID", deviceID))
			rwMu.Lock()
			delete(clientLastHeartbeat, deviceID)
			delete(connectionMap, deviceID)
			rwMu.Unlock()
			return
		}

		switch messageType {
		case websocket.TextMessage:
			messageText := string(message)
			logger.Info("message Text", logger.Any("message", messageText))
			var contentBody ContentBody
			err := json.Unmarshal([]byte(messageText), &contentBody)
			if err != nil {
				logger.Error("ParseTextToJSON error", logger.Err(err), logger.String("origin message is", messageText))
				continue
			}

			switch contentBody.ContentType {
			case Ping:
				clientLastHeartbeat[contentBody.SenderID] = time.Now()
				logger.Info("收到心跳包", logger.Any("sender", contentBody.SenderID), logger.Any("remoteAddr", conn.RemoteAddr().String()))
				mutex.Lock()
				connectionMap[contentBody.SenderID] = ConnectionInfo{
					Conn:              conn,
					RemoteAddr:        remoteAddr,
					LastHeartbeatTime: time.Now(),
				}
				mutex.Unlock()
				message := createMessage(Pong, nil, contentBody.SenderID, contentBody.ClientMsgID)
				err = conn.WriteMessage(messageType, message)
				if err != nil {
					logger.Warn("WriteMessage error", logger.Err(err))
					continue
				}
			case Call, Answer, EndCall:
				if err := handleMessageSend(conn, messageType, contentBody, contentBody.ContentType); err != nil {
					continue
				}
			case SendSMS:
				if err := handleMessageSend(conn, messageType, contentBody, contentBody.ContentType); err != nil {
					continue
				}
			case Idle, Income, NewSMS:
				groupOwnerID, err := getGroupOwnerID(ctx, w, contentBody.SenderID)
				if err != nil {
					continue
				}
				logger.Info("收到话机主动上报的数据", logger.Any("sender", contentBody.SenderID), logger.Any("remoteAddr", conn.RemoteAddr().String()), logger.Any("groupOwnerID", groupOwnerID))
				var newContentBody ContentBody
				if contentBody.ContentType == NewSMS {
					content := make(map[string]interface{})
					content["address"] = contentBody.Content
					content["body"] = contentBody.ReceiveID
					newContentBody = ContentBody{
						ContentType: NewSMS,
						Content:     content,
						SenderID:    contentBody.SenderID,
						ReceiveID:   groupOwnerID,
						ClientMsgID: contentBody.ClientMsgID,
					}
				} else {
					newContentBody = ContentBody{
						ContentType: contentBody.ContentType,
						Content:     contentBody.Content,
						SenderID:    contentBody.SenderID,
						ReceiveID:   groupOwnerID,
						ClientMsgID: contentBody.ClientMsgID,
					}
				}
				sendMessageToUser(conn, newContentBody)
			}
		default:
			logger.Warnf("Unknown message type: %d", messageType)
		}
	}
}

func (w websocketHandler) GetOnlineClients(c *gin.Context) {
	mutex.Lock()
	connectionInfoCopy := make(map[string]ConnectionInfo)
	for key, value := range connectionMap {
		connectionInfoCopy[key] = ConnectionInfo{
			RemoteAddr:        value.RemoteAddr,
			LastHeartbeatTime: value.LastHeartbeatTime,
		}
	}
	mutex.Unlock()

	response.Success(c, gin.H{
		"results": connectionInfoCopy,
		"count":   len(connectionInfoCopy),
	})
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		index := rand.Intn(len(charset))
		sb.WriteByte(charset[index])
	}
	return sb.String()
}

func createMessage(message_type ContentType, content interface{}, receiveID string, msgID string) []byte {
	if msgID == "" {
		msgID = GenerateRandomString(6)
	}
	original := ContentBody{
		ContentType: message_type,
		Content:     content,
		SenderID:    "server",
		ReceiveID:   receiveID,
		ClientMsgID: msgID,
	}
	jsonData, err := json.Marshal(original)
	if err != nil {
		logger.Error("JSON 编码出错", logger.Err(err))
		return []byte("error")
	}
	return jsonData
}

func sendMessageToUser(conn *ws.Conn, contentBody ContentBody) error {
	logger.Infof("正在给 %s 发送 %s", contentBody.ReceiveID, contentBody.Content)
	mutex.Lock()
	info, exists := connectionMap[contentBody.ReceiveID]
	mutex.Unlock()
	if !exists {
		logger.Errorf("用户未找到", logger.Any("设备ID", contentBody.ReceiveID))
		return errors.New("用户未找到")
	}
	message, err := json.Marshal(contentBody)
	if err != nil {
		logger.Error("JSON 编码出错", logger.Err(err))
		return errors.New("JSON 编码出错")
	}

	if info.Conn != nil {
		err := info.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.Infof("给用户 %s 发送消息失败: %v", contentBody.ReceiveID, err)
			return errors.New("给用户发送消息失败")
		}
		logger.Infof("给用户 %s 发送消息成功: %s", contentBody.ReceiveID, string(message))
		message := createMessage(SendSuccess, contentBody.ContentType, contentBody.ReceiveID, contentBody.ClientMsgID)
		return conn.WriteMessage(websocket.TextMessage, message)
	}
	logger.Infof("用户不存在", logger.Any("设备ID", contentBody.ReceiveID))
	return errors.New("用户不存在")
}
