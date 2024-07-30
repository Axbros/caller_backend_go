package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		clientsRouter(group, handler.NewClientsHandler())
	})
}

func clientsRouter(group *gin.RouterGroup, h handler.ClientsHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/clients", h.Create)
	group.DELETE("/clients/:id", h.DeleteByID)
	group.PUT("/clients/:id", h.UpdateByID)
	group.GET("/clients/:id", h.GetByID)
	group.POST("/clients/list", h.List)

	group.POST("/clients/delete/ids", h.DeleteByIDs)
	group.POST("/clients/condition", h.GetByCondition)
	group.POST("/clients/list/ids", h.ListByIDs)
	group.GET("/clients/list", h.ListByLastID)
}
