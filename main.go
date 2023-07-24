package main

import (
	"github.com/spf13/viper"
	"imessage/router"
	"imessage/utils"
)

func main() {
	// 初始化
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	// 路由
	r := router.Router()
	err := r.Run(viper.GetString("port.server"))
	if err != nil {
		return
	}
}
