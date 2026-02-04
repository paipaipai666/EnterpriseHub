package initializers

import (
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		Log.Fatal("DATABASE_URL 环境变量未设置")
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 添加日志以便调试
	})

	if err != nil {
		Log.Fatal("连接数据库失败: %v", zap.Error(err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		Log.Fatal("获取底层数据库连接失败: %v", zap.Error(err))
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		Log.Fatal("数据库 ping 失败: %v", zap.Error(err))
	}

	Log.Info("MySQL数据库连接成功")
}
