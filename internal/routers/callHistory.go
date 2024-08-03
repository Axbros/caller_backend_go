package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		callHistoryRouter(group, handler.NewCallHistoryHandler())
	})
}

func callHistoryRouter(group *gin.RouterGroup, h handler.CallHistoryHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/callHistory", h.Create)
	group.DELETE("/callHistory/:id", h.DeleteByID)
	group.PUT("/callHistory/:id", h.UpdateByID)
	group.GET("/callHistory/:id", h.GetByID)
	group.POST("/callHistory/list", h.List)

	group.POST("/callHistory/delete/ids", h.DeleteByIDs)
	group.POST("/callHistory/condition", h.GetByCondition)
	group.POST("/callHistory/list/ids", h.ListByIDs)
	group.GET("/callHistory/list", h.ListByLastID)
}
