package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitSrvConn() { //这个采取了负载均衡的方式连接用户服务，可以实现动态的负载均衡,重点在_ "github.com/mbobakov/grpc-consul-resolver"
	// 从全局配置中获取Consul信息
	//_ "github.com/mbobakov/grpc-consul-resolver"虽然没有在代码中显式调用，
	//但是它有init()函数，会在程序启动时自动执行，从而完成初始化工作
	consulInfo := global.ServerConfig.ConsulInfo

	// 使用gRPC Dial函数连接到Consul注册的用户服务
	userConn, err := grpc.Dial(
		// 格式化Consul服务地址
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=srv", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		// 使用不安全的证书创建传输凭证（因为是本地开发，不使用TLS）
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 设置默认的服务配置，使用轮询负载均衡策略
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	// 如果连接失败，记录错误并终止程序
	if err != nil {
		zap.S().Fatal("InitSrvConn用户服务连接失败")
	}

	// 创建用户服务客户端
	userSrvClient := proto.NewUserClient(userConn)

	// 将客户端保存到全局变量中
	global.UserSrvClient = userSrvClient

}

func InitSrvConn2() {
	//从注册中心获取用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))

	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	//这个程序没有办法从cousil去读取配置
	//userSrvHost = "192.168.43.127"
	//userSrvPort = 50051
	if userSrvHost == "" {

		zap.S().Fatal("[InitSrvConn]链接用户服务失败")
		return
	}

	//拨号连接用户grpc服务 跨域问题，后端和前端都可以解决
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("GetUserList 连接用户服务失败",
			"msg", err.Error(),
		)
	}
	//1,已经事先创立好了链接，这样后续就不用进行tcp的三次握手了，直接就可以进行grpc的通信了
	//2，一个链接可以支持多个groutine的请求，要使用链接池
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
