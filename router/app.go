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

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 静态资源
	r.Static("/asset", "/Users/luliang/GoLand/imessage/asset/")
	r.LoadHTMLGlob("/Users/luliang/GoLand/imessage/views/**/*")

	// 首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)        // 用户注册
	r.GET("/toChat", service.ToChat)                // 聊天页面
	r.POST("/contact/addFriend", service.AddFriend) // 添加好友页面
	r.POST("/searchFriends", service.SearchFriends) // 返回好友列表
	// 群
	r.POST("/contact/createCommunity", service.CreateCommunity) //创建群
	r.POST("/contact/loadcommunity", service.LoadCommunity)     //群列表
	r.POST("/contact/joinGroup", service.JoinGroups)            // 添加群聊
	// 用户模块
	r.POST("/user/createUser", service.CreateUser) // 增加用户
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("/user/login", service.FindUserByNameAndPwd)
	r.POST("/user/find", service.FindByID)

	// 发送消息
	r.GET("/user/sendMsg", service.SendMsg) // websocket 测试
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.GET("/chat", service.Chat)
	r.POST("/user/redisMsg", service.RedisMsg)

	return r
}
