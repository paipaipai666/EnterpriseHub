package main

import (
	"fmt"
	"log"

	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

func main() {
	if initializers.DB == nil {
		log.Fatal("数据库连接未初始化")
	}

	if err := initializers.DB.AutoMigrate(&model.User{}); err != nil {
		log.Printf("迁移 User 表失败: %v", err)
		return
	}

	fmt.Println("数据库迁移完成")
}
