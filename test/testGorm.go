package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"imessage/models"
)

// 这里就是将模板导入数据库的一个简单案例
// 前提是创建好了一个数据库,数据库的名字暂定义为 "ginchat"

var (
	DB *gorm.DB
)

func initMySQL() (err error) {
	// 连接数据库
	dsn := "root:791975457@qq.com@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 判断是否连通
	return DB.DB().Ping()
}
func main() {
	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	// 模型绑定
	DB.AutoMigrate(&models.UserBasic{}) // todos
	defer func(DB *gorm.DB) {
		err := DB.Close()
		if err != nil {
			panic(err)
		}
	}(DB)
	var user1 models.UserBasic

	user1.Name = "Lucy"
	user1.PassWord = "123456"
	user1.Phone = "110120119"
	user1.Email = "123@qq.com"
	user1.Identity = "user1"
	user1.ClientIP = "127.1.1.1"
	user1.ClientPort = "9000"
	user1.LoginTime = 1
	user1.HeartBeatTime = 1
	user1.LoginOutTime = 1
	user1.IsLogOut = false
	user1.DeviceInfo = "iOS"

	result := DB.Create(&user1)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}
