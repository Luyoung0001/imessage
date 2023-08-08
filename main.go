package main

import (
	"github.com/spf13/viper"
	"server_imessage/models"
	"server_imessage/router"
	"server_imessage/utils"
	"time"
)

func main() {
	// 初始化
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	InitTimer()
	// 路由
	r := router.Router()
	err := r.Run(viper.GetString("port.server"))
	if err != nil {
		return
	}
}

// 定时清理死掉的连接

func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second, time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}
