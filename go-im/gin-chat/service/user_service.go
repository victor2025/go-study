package service

import (
	"gin-chat/models"
	"gin-chat/utils"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// getUserList
// @Summary 获取用户列表
// @Tags 用户模块
// @Schemes
// @Description get user list
// @Success 200 {string} []*UserBasic
// @Router /user/list [get]
func GetUserList(c *gin.Context) {
	db, data := models.GetUserList()
	c.JSON(200, gin.H{
		"success": db.RowsAffected,
		"message": data,
	})
}

// findUser
// @Summary 查找指定用户
// @Tags 用户模块
// @Param name query string false "username"
// @Param phone query string false "phone"
// @Param email query string false "email"
// @Schemes
// @Description find user
// @Success 200 {string} user
// @Router /user/find [get]
func FindUser(c *gin.Context) {
	user := &models.UserBasic{
		Name:  c.Query("name"),
		Phone: c.Query("phone"),
		Email: c.Query("email"),
	}
	db := models.FindUser(user)
	c.JSON(200, gin.H{
		"success": db.RowsAffected,
		"message": user,
	})
}

// createUser
// @Summary 新建用户
// @Tags 用户模块
// @Param name query string false "用户名"
// @Param password query string false "密码"
// @Param repassword query string false "重复输入密码"
// @Schemes
// @Description create new user
// @Success 200 {string} success
// @Router /user/create [get]
func CreateUser(c *gin.Context) {
	user := &models.UserBasic{}
	user.Name = c.Query("name")

	// 查找是否有重名用户
	db := models.FindUser(user)
	if db.RowsAffected > 0 {
		c.JSON(-1, gin.H{
			"success": 0,
			"message": "用户名已存在",
		})
		return
	}

	// 比对密码
	password := c.Query("password")
	repassword := c.Query("repassword")
	if password != repassword {
		c.JSON(-1, gin.H{
			"success": 0,
			"message": "两次密码不同",
		})
		return
	}
	user.Password = password
	success := models.CreateUser(user).RowsAffected

	c.JSON(200, gin.H{
		"success": success,
		"message": "注册成功",
	})
}

// deleteUser
// @Summary 删除用户
// @Tags 用户模块
// @Param id query string false "id"
// @Schemes
// @Description delete user
// @Success 200 {string} success
// @Router /user/delete [get]
func DeleteUser(c *gin.Context) {
	user := &models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)

	success := models.DeleteUser(user).RowsAffected

	c.JSON(200, gin.H{
		"success": success,
		"message": "",
	})
}

// updateUser
// @Summary 修改用户
// @Tags 用户模块
// @Param id formData string false "id"
// @Param name formData string false "username"
// @Param password formData string false "password"
// @Param phone formData string false "phone"
// @Param email formData string false "email"
// @Schemes
// @Description update user
// @Success 200 {string} success
// @Router /user/update [post]
func UpdateUser(c *gin.Context) {
	user := &models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	// 验证数据并响应
	_, err := govalidator.ValidateStruct(user)
	utils.HandleError(err, func() {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "数据格式不合法",
		})
	}, func() {
		success := models.UpdateUser(user).RowsAffected
		c.JSON(200, gin.H{
			"success": success,
			"message": "修改成功",
		})
	})

}
