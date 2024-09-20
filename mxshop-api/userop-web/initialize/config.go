package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/userop-web/global"
)

// GetEnvInfo 获取环境变量信息
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()      // AutomaticEnv 自动读取环境变量
	return viper.GetBool(env) // GetBool 获取与键关联的值作为布尔值
}

// InitConfig 初始化配置
func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG") // 检查是否启用了调试模式

	var configFileName string
	configFileNamePrefix := "config"

	// 根据调试模式选择配置文件
	if debug {
		configFileName = fmt.Sprintf("userop-web/%s-debug.yaml", configFileNamePrefix)
	} else {
		configFileName = fmt.Sprintf("userop-web/%s-pro.yaml", configFileNamePrefix)
	}

	v := viper.New()                // 创建一个新的 viper 实例
	v.SetConfigFile(configFileName) // 设置要读取的配置文件

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		panic(err) // 如果读取配置文件出错，触发 panic
	}

	// 将配置反序列化为全局 ServerConfig 结构体
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err) // 如果反序列化出错，触发 panic
	}

	// 打印配置信息到控制台
	zap.S().Infof("nacos配置信息：%v", global.NacosConfig) // 使用 Zap 记录日志输出配置信息
	//从nacos中获取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败:%s", err.Error())
	}
	fmt.Println(&global.ServerConfig)
	//fmt.Printf("%v\n", v.Get("name"))             // 打印特定配置项到控制台

	//v.WatchConfig()                           // 监听配置文件变化
	//v.OnConfigChange(func(e fsnotify.Event) { // 当配置文件发生变化时执行回调函数
	//	zap.S().Infof("配置文件变更:%v", e.Name)            // 使用 Zap 记录日志输出配置文件变更信息
	//	_ = v.ReadInConfig()                          // 重新读取配置文件
	//	_ = v.Unmarshal(&global.ServerConfig)         // 重新反序列化配置到全局 ServerConfig 结构体
	//	zap.S().Infof("配置信息：%v", global.ServerConfig) // 使用 Zap 记录日志输出更新后的配置信息
	//})
}
