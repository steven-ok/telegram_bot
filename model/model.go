package model

import (
	"tg_bots/common"
	"tg_bots/model/topic"
)

func Init() {
	// 执行模型迁移
	common.DB.AutoMigrate(&topic.Topic{}, &topic.TopicReply{}, &topic.TopicAdvocate{})
}
