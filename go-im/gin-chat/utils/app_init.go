package utils

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{})
	HandleError(err)
	_ = DB
}
