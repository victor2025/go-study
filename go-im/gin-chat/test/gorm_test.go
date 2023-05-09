package test

import (
	"gin-chat/models"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestBasicFunctions(t *testing.T) {
	dbPath := "victor2022:1045899571@tcp(172.30.1.2)/gin_chat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	t.Logf("database connected")

	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})

	// Create
	user := &models.UserBasic{
		Name: "victor2022",
	}
	db.Create(user)

	// Read
	db.First(&user, 1) // 根据整型主键查找
	t.Logf("get user: %v", user)
	// db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 000000
	db.Model(&user).Update("Password", 000000)
	// Update - 更新多个字段
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	// if isDel := db.Unscoped().Delete(&user); isDel.RowsAffected == 0 {
	// 	t.Error("delete operation failed")
	// }

}
