package router

import (
	"gin-chat/docs"
	"gin-chat/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	swaggerRouter(r)
	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.GET("/user/create", service.CreateUser)
	r.GET("/user/delete", service.DeleteUser)
	r.POST("/user/update", service.UpdateUser)
	r.GET("/user/find", service.FindUser)
	r.POST("/user/login", service.Login)
	return r
}

func swaggerRouter(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
