package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_sevs/goods_srv/global"
	"mxshop_sevs/goods_srv/handler"
	"mxshop_sevs/goods_srv/initialize"
	"mxshop_sevs/goods_srv/model/proto"
	"mxshop_sevs/goods_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitEs()
	zap.S().Info(global.ServerConfig)
	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip:", *IP, "port:", *Port)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check
	//1.如何启动两个服务
	//2.即使通过终端启动两个服务，注册到consul中，consul中也只会有一个服务实例，因为名字相同，所以只能有一个实例
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis) //这是一个阻塞的函数，直到server.Stop()被调用才会返回，因此，需要在另一个goroutine中调用server.Stop()
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()
	//优雅的退出
	// 创建一个通道用于接收操作系统的退出信号
	quit := make(chan os.Signal)
	// 监听操作系统的中断信号和终止信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞，直到接收到退出信号
	<-quit
	// 尝试注销服务，如果失败则记录日志
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	// 记录日志，表示服务已关闭
	zap.S().Info("服务关闭")

}
