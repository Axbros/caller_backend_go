package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		unanswerdCallRouter(group, handler.NewUnanswerdCallHandler())
	})
}

func unanswerdCallRouter(group *gin.RouterGroup, h handler.UnanswerdCallHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/unanswerdCall", h.Create)
	group.DELETE("/unanswerdCall/:id", h.DeleteByID)
	group.PUT("/unanswerdCall/:id", h.UpdateByID)
	group.GET("/unanswerdCall/:id", h.GetByID)
	group.POST("/unanswerdCall/list", h.List)

	group.POST("/unanswerdCall/delete/ids", h.DeleteByIDs)
	group.POST("/unanswerdCall/condition", h.GetByCondition)
	group.POST("/unanswerdCall/list/ids", h.ListByIDs)
	group.GET("/unanswerdCall/list", h.ListByLastID)
}
