package service

import (
	"github.com/gin-gonic/gin"
	"imessage/models"
)

func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})

}
