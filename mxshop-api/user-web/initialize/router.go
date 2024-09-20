package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"
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
	Router.Use(middlewares.Cors()) // 跨域中间件
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup) // 注册用户相关路由，目前路由就是/v1/user/list,路由组的路由是可以叠加的
	router.InitBaseRouter(ApiGroup) //注册图形验证码相关路由
	return Router
}
