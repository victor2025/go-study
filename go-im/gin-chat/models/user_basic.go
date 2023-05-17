package models

import (
	"fmt"
	"gin-chat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model    // 基本字段，其中引入了软删除字段
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9})"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string    // 密码盐
	LoginTime     time.Time // 上线时间
	HeartbeatTime time.Time // 心跳时间
	LogoutTime    time.Time // 下线时间
	IsLogout      bool
	DeviceInfo    string // 设备信息
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() (*gorm.DB, []*UserBasic) {
	data := make([]*UserBasic, 100)
	db := utils.DB.Find(&data)
	for _, v := range data {
		fmt.Printf("v: %v\n", v)
	}
	return db, data
}

func FindUser(user *UserBasic) *gorm.DB {
	return utils.DB.Where(user).First(user)
}

func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}

func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(user)
}
