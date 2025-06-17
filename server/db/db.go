package db

import (
	"log"

	// "github.com/erjean0310/examgen/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库初始化失败: ", err.Error())
	}
	// 自动迁移数据库
	// DB.AutoMigrate(&models.Question{})
}
