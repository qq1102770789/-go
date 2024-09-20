package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop-api/order-web/global"
	"mxshop-api/order-web/proto"
	"mxshop-api/order-web/utils/otgrpc"
)

func InitSrvConn() { //这个采取了负载均衡的方式连接用户服务，可以实现动态的负载均衡,重点在_ "github.com/mbobakov/grpc-consul-resolver"
	// 从全局配置中获取Consul信息
	//_ "github.com/mbobakov/grpc-consul-resolver"虽然没有在代码中显式调用，
	//但是它有init()函数，会在程序启动时自动执行，从而完成初始化工作
	consulInfo := global.ServerConfig.ConsulInfo

	// 使用gRPC Dial函数连接到Consul注册的用户服务
	goodsConn, err := grpc.Dial(
		// 格式化Consul服务地址
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=srv", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		// 使用不安全的证书创建传输凭证（因为是本地开发，不使用TLS）
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 设置默认的服务配置，使用轮询负载均衡策略
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	// 如果连接失败，记录错误并终止程序
	if err != nil {
		zap.S().Fatal("InitSrvConn商品服务连接失败")
	}

	// 创建用户服务客户端
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	// 将客户端保存到全局变量中

	orderConn, err := grpc.Dial(
		// 格式化Consul服务地址
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=srv", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		// 使用不安全的证书创建传输凭证（因为是本地开发，不使用TLS）
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 设置默认的服务配置，使用轮询负载均衡策略
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	// 如果连接失败，记录错误并终止程序
	if err != nil {
		zap.S().Fatal("InitSrvConn订单服务连接失败")
	}

	// 创建用户服务客户端
	global.OrderSrvClient = proto.NewOrderClient(orderConn)
	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}

	global.InventorySrvClient = proto.NewInventoryClient(invConn)

}
