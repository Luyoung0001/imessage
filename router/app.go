package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"imessage/docs"
	"imessage/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList) // 查
	r.GET("/user/createUser", service.CreateUser)   // 增
	r.GET("/user/deleteUser", service.DeleteUser)   // 删
	r.POST("/user/updateUser", service.UpdateUser)  // 改

	return r
}
