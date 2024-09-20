package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/userop-web/global"
	"mxshop-api/userop-web/initialize"
	"mxshop-api/userop-web/utils"
	"mxshop-api/userop-web/utils/register/consul"
	myvalidator "mxshop-api/userop-web/validator"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	initialize.InitLogger()                            //初始化日志
	initialize.InitConfig()                            ////初始化配置
	if err := initialize.InitTrans("zh"); err != nil { //初始化翻译
		panic(err)
	}
	Router := initialize.Routers() //需要返回路由，是因为我们要在main函数里启动

	initialize.InitSrvConn()

	viper.AutomaticEnv() // AutomaticEnv 自动读取环境变量
	//如果是本地开发环境，则固定端口，否则随机获取一个空闲端口
	debug := viper.GetBool("MXSHOP_DEBUG") // GetBool 获取与键关联的值作为布尔值
	if !debug {
		port, err := utils.GetFreePort()

		if err == nil {
			global.ServerConfig.Port = port
		}

	}
	//注册自定义验证器,配置是通过标签来实现的，比如：mobile:"required,mobile"

	//s()是获取一个全局的sugar，可以让我们自己设置一个全局的logger，同时返回的sugar是安全的，因为他有锁机制，可以保证协程安全
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 如果可以获取到验证引擎，就将自定义的手机号码验证函数注册到验证引擎中
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		// 然后注册手机号码验证的翻译信息
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			// 添加手机号码验证错误的翻译信息，{0} 会被替换为字段名称
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // 详细信息请参考 universal-translator
		}, func(ut ut.Translator, fe validator.FieldError) string {
			// 返回手机号码验证错误的翻译信息，使用字段名称替换 {0}
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceStr := uuid.NewV4()
	serviceId := serviceStr.String()
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId) //注册服务到consul
	if err != nil {
		zap.S().Panic("注册服务到consul失败:", err.Error())
	}
	zap.S().Debugf("启动服务器，监听端口：%d", global.ServerConfig.Port) //日志记录.S()是获取全局日志对象

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败:", err.Error)
		}
	}()
	//接受终止信号，优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.Deregister(serviceId); err != nil {
		zap.S().Panic("注销服务失败:", err.Error())
	} else {
		zap.S().Info("注销服务成功")
	} //注销服务

}
