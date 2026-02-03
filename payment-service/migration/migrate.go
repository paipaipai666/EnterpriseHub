package main

import (
	"fmt"
	"log"

	"github.com/paipaipai666/EnterpriseHub/payment-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/domain"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

func main() {
	if initializers.DB == nil {
		log.Fatal("数据库连接未初始化")
	}

	if err := initializers.DB.AutoMigrate(&domain.Payment{}); err != nil {
		log.Printf("迁移 Payment 表失败: %v", err)
		return
	}

	fmt.Println("数据库迁移完成")
}
