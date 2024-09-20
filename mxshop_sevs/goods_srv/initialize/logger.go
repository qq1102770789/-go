package initialize

import "go.uber.org/zap"

// InitLogger 初始化应用程序的日志记录器。
func InitLogger() {
	// 使用 zap.NewDevelopment() 创建一个新的开发环境日志记录器实例。
	logger, _ := zap.NewDevelopment()
	// 使用 zap.ReplaceGlobals() 将全局日志记录器替换为新创建的日志记录器。
	zap.ReplaceGlobals(logger)
}
