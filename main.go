package main

import (
	"imessage/router"
	"imessage/utils"
)

func main() {
	// 初始化
	utils.InitConfig()
	utils.InitMySQL()
	// 路由
	r := router.Router()
	err := r.Run(":9000")
	if err != nil {
		return
	}
}
