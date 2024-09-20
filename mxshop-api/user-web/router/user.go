package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

// InitUserRouter 初始化用户路由。
func InitUserRouter(Router *gin.RouterGroup) {
	// 在传入的 RouterGroup 上创建一个名为 "user" 的子路由组。
	UserRouter := Router.Group("user") //加入JWT验证中间件
	{
		// 在 "user" 路由组下创建 GET 请求处理函数为 api.GetUserList 的路由 "list"。
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)

		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}
}
