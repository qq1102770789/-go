package middlewares

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mxshop-api/userop-web/global"
	"mxshop-api/userop-web/models"
	"net/http"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			var errMsg string
			if err == TokenExpired {
				errMsg = "授权已过期"
			} else if err == TokenMalformed {
				errMsg = "Token格式错误"
			} else if err == TokenNotValidYet {
				errMsg = "Token未生效"
			} else if err == TokenInvalid {
				errMsg = "无效的Token"
			} else {
				errMsg = "无法处理该Token"
			}

			// 调试信息
			fmt.Println("JWTAuth error:", errMsg, "Token:", token, "Error:", err)

			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"msg":   "jwt未登陆",
				"error": errMsg,
				"token": token,
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.SigningKey)
	if err != nil {
		// 调试信息
		fmt.Println("CreateToken error:", err)
	}
	return tokenString, err
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 调试信息
	fmt.Println("Parsing token:", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// 调试信息
			fmt.Println("ValidationError occurred:", ve)

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
