package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_sevs/order_srv/global"
	"mxshop_sevs/order_srv/model/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	goodsConn, err := grpc.Dial(
		// 格式化Consul服务地址
		fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=srv", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		// 使用不安全的证书创建传输凭证（因为是本地开发，不使用TLS）
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 设置默认的服务配置，使用轮询负载均衡策略
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
	}

	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	//初始化库存服务连接
	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
	}

	global.InventorySrvClient = proto.NewInventoryClient(invConn)
}
