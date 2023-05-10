package utils

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func InitApp() {
	initConfig()
	initMySQL()
}

func initConfig() {
	// 读取配置文件
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	HandleError(err, func() {
		os.Exit(1)
	})
	keys := viper.AllKeys()
	for _, v := range keys {
		log.Printf("Get config \"%v\": %v", v, viper.Get(v))
	}
}

func initMySQL() {
	mysqlLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")),
		&gorm.Config{Logger: mysqlLogger})
	HandleError(err)
	_ = DB
}
