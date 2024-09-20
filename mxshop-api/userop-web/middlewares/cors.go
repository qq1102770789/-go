package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 中间件用于处理跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求方法
		method := c.Request.Method
		// 设置允许跨域的域名，* 表示允许任意域名访问，也可以指定特定域名
		c.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,x-token")
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST ,GET, OPTIONS, PUT, DELETE, PATCH")
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 设置是否允许发送 Cookie，true 表示允许
		c.Header("Access-Control-Allow-Credentials", "true")
		// 如果请求方法是 OPTIONS，则返回 No Content 状态码
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

//
