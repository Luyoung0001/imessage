package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DBGinChat *gorm.DB
)

func initMySQLGinChat() (err error) {
	// 连接数据库
	dsn := "root:791975457@qqCom@tcp(rm-cn-pe33bz99b000emxo.rwlb.rds.aliyuncs.com:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
	DBGinChat, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 判断是否连通
	fmt.Println("PONG!")
	return DBGinChat.DB().Ping()
}
func main() {
	// 连接数据库
	err := initMySQLGinChat()
	if err != nil {
		panic(err)
		return
	}

}
