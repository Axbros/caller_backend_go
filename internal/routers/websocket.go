package routers

import (
	"caller/internal/handler"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/ws"
)

func init() {
	websocketRouterFns = append(websocketRouterFns, func(group *gin.RouterGroup) {
		websocketRouter(group, handler.NewWebsocketHandler())
	})
}

func websocketRouter(group *gin.RouterGroup, h handler.WebsocketHandler) {

	group.GET("/", func(c *gin.Context) {
		{
			s := ws.NewServer(c.Writer, c.Request, h.LoopReceiveMessage) // default setting
			err := s.Run(context.Background())
			if err != nil {
				log.Println("webSocket server error:", err)

			}
		}
	})
	group.GET("/online", h.GetOnlineClients)
}
