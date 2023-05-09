package models

import (
	"fmt"
	"gin-chat/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model    // 基本字段，其中引入了软删除字段
	Name          string
	Password      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64 // 上线时间
	HeartbeatTime uint64 // 心跳时间
	LogoutTime    uint64 // 下线时间
	IsLogout      bool
	DeviceInfo    string // 设备信息
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() *[]*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Printf("v: %v\n", v)
	}
	return &data
}
