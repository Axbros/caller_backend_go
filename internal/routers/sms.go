package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		smsRouter(group, handler.NewSmsHandler())
	})
}

func smsRouter(group *gin.RouterGroup, h handler.SmsHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/sms", h.Create)
	group.DELETE("/sms/:id", h.DeleteByID)
	group.PUT("/sms/:id", h.UpdateByID)
	group.GET("/sms/:id", h.GetByID)
	group.POST("/sms/list", h.List)

	group.POST("/sms/delete/ids", h.DeleteByIDs)
	group.POST("/sms/condition", h.GetByCondition)
	group.POST("/sms/list/ids", h.ListByIDs)
	group.GET("/sms/list", h.ListByLastID)
}
