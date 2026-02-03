package main

import (
	"fmt"
	"log"

	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

func main() {
	if initializers.DB == nil {
		log.Fatal("数据库连接未初始化")
	}

	if err := initializers.DB.AutoMigrate(&domain.Order{}); err != nil {
		log.Printf("迁移 Order 表失败: %v", err)
		return
	}

	fmt.Println("数据库迁移完成")
}
