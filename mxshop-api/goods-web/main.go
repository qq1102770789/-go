package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialize"
	"mxshop-api/goods-web/utils"
	"mxshop-api/goods-web/utils/register/consul"
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
	initialize.InitSentinel()

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
