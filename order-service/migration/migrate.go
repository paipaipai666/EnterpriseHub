package main

import (
	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
	"go.uber.org/zap"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("order_service")
	initializers.ConnectToDatabase()
}

func main() {
	if initializers.DB == nil {
		initializers.Log.Fatal("数据库连接未初始化")
	}

	if err := initializers.DB.AutoMigrate(&domain.Order{}); err != nil {
		initializers.Log.Fatal("迁移 Order 表失败: %v", zap.Error(err))
		return
	}

	initializers.Log.Info("数据库迁移完成")
}
