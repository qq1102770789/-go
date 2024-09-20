package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/order-web/middlewares"
	"mxshop-api/order-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())
	//添加链路追踪
	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup) // 注册用户相关路由，目前路由就是/v1/user/list,
	router.InitShopCartRouter(ApiGroup)

	return Router
}
