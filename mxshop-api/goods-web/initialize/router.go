package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/middlewares"
	"mxshop-api/goods-web/router"
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
	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup) // 注册用户相关路由，目前路由就是/v1/user/list,路由组的路由是可以叠加的
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)
	return Router
}
