package topic

import (
	"github.com/jinzhu/gorm"
	"math/rand"
	"tg_bots/common"
)

type TopicAdvocate struct {
	gorm.Model
	OwnerId int    `json:"owner_id" gorm:"desc:'拥护者ID'"`
	Owner   string `json:"owner" gorm:"desc:'拥护者姓名'"`
	ChatId  int64  `json:"chat_id" gorm:"desc:'拥护会话'"`
}

func TopicAdvocatesByChatId(chatId int64) []TopicAdvocate {
	var list []TopicAdvocate
	common.DB.Model(&TopicAdvocate{}).Where("chat_id = ?", chatId).Find(&list)

	for i := len(list) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		list[i], list[num] = list[num], list[i]
	}
	return list
}
