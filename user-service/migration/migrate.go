package main

import (
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
	"go.uber.org/zap"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("user_service")
	initializers.ConnectToDatabase()
}

func main() {
	if initializers.DB == nil {
		initializers.Log.Fatal("数据库连接未初始化")
	}

	if err := initializers.DB.AutoMigrate(&model.User{}); err != nil {
		initializers.Log.Fatal("迁移 User 表失败: %v", zap.Error(err))
		return
	}

	initializers.Log.Info("数据库迁移完成")
}
