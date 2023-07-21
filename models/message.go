package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FromId   string
	TargetId string
	Type     string // 群聊,私聊,广播等
	Name     string
	Media    int    // 文字,图片,音频等
	Content  string // 消息内容
	Pic      string // 图片
	Url      string // 链接
	Desc     string // 描述
	Amount   int
}

func (table *Message) TableName() string {
	return "message"
}
