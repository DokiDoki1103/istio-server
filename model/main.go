package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"istio-server/kubernetes"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB(k8sClient *kubernetes.K8SClient) (err error) {

	var db *gorm.DB
	dsn := os.Getenv("SQL_DSN")
	if dsn == "" {
		pwd, err := k8sClient.GetMySQLPassword()
		if err != nil {
			log.Fatalf("Failed to get MySQL password: %s", err)
		}
		dsn = fmt.Sprintf("root:%s@tcp(rbd-db-rw.rbd-system.svc.cluster.local:3306)/console?charset=utf8mb4&parseTime=True&loc=Local", pwd)
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Info), // 启用日志模式，可以选择 logger.Silent, logger.Error, logger.Warn, logger.Info
	})
	if err == nil {
		DB = db
	}
	return err
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
