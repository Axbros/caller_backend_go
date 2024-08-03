package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		groupClientRouter(group, handler.NewGroupClientHandler())
	})
}

func groupClientRouter(group *gin.RouterGroup, h handler.GroupClientHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/groupClient", h.Create)
	group.DELETE("/groupClient/:id", h.DeleteByID)
	group.PUT("/groupClient/:id", h.UpdateByID)
	group.GET("/groupClient/:id", h.GetByID)
	group.POST("/groupClient/list", h.List)

	group.POST("/groupClient/delete/ids", h.DeleteByIDs)
	group.POST("/groupClient/condition", h.GetByCondition)
	group.POST("/groupClient/list/ids", h.ListByIDs)
	group.GET("/groupClient/list", h.ListByLastID)
}
