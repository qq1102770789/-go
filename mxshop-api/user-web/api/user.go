package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": +e.Code(),
				})
			}
			return

		}
	}
	//获取错误状态码并转化为http状态码

}
func HandleValidatorError(c *gin.Context, err error) {

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return

}
func GetUserList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser.ID)
	//生成grpc的client并调用接口

	pn := ctx.DefaultQuery("pn", "0") //是从url中获取参数
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "0")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("GetUserList 获取用户列表失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Mobile:   value.Mobile,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
		}

		result = append(result, user)

	}
	ctx.JSON(http.StatusOK, result)

}
func PassWordLogin(c *gin.Context) {
	// 初始化密码登录表单结构体
	passwordLoginForm := forms.PassWordLoginForm{}
	// 绑定请求参数到表单结构体
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err) // 处理验证器错误
		return
	}
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}
	// 拨号连接用户 gRPC 服务

	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		// 处理 gRPC 调用错误
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				// 用户不存在的情况
				c.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				// 其他错误情况
				c.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}
	} else {
		// 只是查询到了用户，但是没有密码，不能登录
		if passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PassWordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			// 处理密码校验失败的情况
			c.JSON(http.StatusInternalServerError, gin.H{
				"password": "登录失败",
			})
		} else {
			// 处理密码校验成功的情况
			if passRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //签名的失效时间,30天
						Issuer:    "imooc",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "登录失败",
				})
			}
		}
	}
}

func Register(c *gin.Context) {
	// 初始化注册表单结构体
	registerForm := forms.RegisterForm{}
	// 绑定请求参数到表单结构体
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err) // 处理验证器错误
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}

	// 生成 gRPC 的 client 并调用接口

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("新建用户失败：%s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //签名的失效时间,30天
			Issuer:    "imooc",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})

}
