package routers

import (
	"github.com/gin-gonic/gin"

	"caller/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		distributionRouter(group, handler.NewDistributionHandler())
	})
}

func distributionRouter(group *gin.RouterGroup, h handler.DistributionHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/distribution", h.Create)
	group.DELETE("/distribution/:id", h.DeleteByID)
	group.PUT("/distribution/:id", h.UpdateByID)
	group.GET("/distribution/:id", h.GetByID)
	group.POST("/distribution/list", h.List)

	group.POST("/distribution/delete/ids", h.DeleteByIDs)
	group.POST("/distribution/condition", h.GetByCondition)
	group.POST("/distribution/list/ids", h.ListByIDs)
	group.GET("/distribution/list", h.ListByLastID)
}
