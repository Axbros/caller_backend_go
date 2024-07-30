package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		groupCallRouter(group, handler.NewGroupCallHandler())
	})
}

func groupCallRouter(group *gin.RouterGroup, h handler.GroupCallHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/groupCall", h.Create)
	group.DELETE("/groupCall/:id", h.DeleteByID)
	group.PUT("/groupCall/:id", h.UpdateByID)
	group.GET("/groupCall/:id", h.GetByID)
	group.POST("/groupCall/list", h.List)

	group.POST("/groupCall/delete/ids", h.DeleteByIDs)
	group.POST("/groupCall/condition", h.GetByCondition)
	group.POST("/groupCall/list/ids", h.ListByIDs)
	group.GET("/groupCall/list", h.ListByLastID)
}
