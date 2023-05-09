package service

import "github.com/gin-gonic/gin"

// @Schemes
// @Description index
// @Success 200 {string} index
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "index",
	})
}
