package main

import (
	"gin-chat/router"
	"gin-chat/utils"
)

func main() {
	utils.InitApp()

	r := router.Router()
	r.Run(":8081")
}
