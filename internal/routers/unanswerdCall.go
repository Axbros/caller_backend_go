package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		callLogRouter(group, handler.NewUnanswerdCallHandler())
	})
}

func callLogRouter(group *gin.RouterGroup, h handler.UnanswerdCallHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/callLog", h.Create)
	group.POST("/callLog/fromDevice/:type", h.MultipleCreate)
	group.DELETE("/callLog/:id", h.DeleteByID)
	group.PUT("/callLog/:id", h.UpdateByID)
	group.GET("/callLog/:id", h.GetByID)
	group.POST("/callLog/list", h.List)

	group.POST("/callLog/delete/ids", h.DeleteByIDs)
	group.POST("/callLog/condition", h.GetByCondition)
	group.POST("/callLog/list/ids", h.ListByIDs)
	group.GET("/callLog/list", h.ListByLastID)
	group.GET("/callLog/byUserId/:id", h.GetByUserID)
}
