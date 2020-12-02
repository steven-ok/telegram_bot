package topic

import (
	"github.com/jinzhu/gorm"
	"sync"
	"tg_bots/common"
)

var Lock sync.Mutex
var Lock2 sync.Mutex

type Topic struct {
	gorm.Model
	Owner         string `json:"owner" gorm:"size:255;desc:发起人姓名"`
	OwnerId       int    `json:"owner_id" gorm:"type:int;desc:发起人ID"`
	Title         string `json:"title" gorm:"desc:接龙标题"`
	ChatId        int64  `json:"chat_id" gorm:"desc:所在会话"`
	ReplayContent string `json:"replay_content" gorm:"type:text;desc:回复内容"`
	LastMessageId int    `json:"last_message_id" gorm:"desc:最后一次消息ID"`
}

func NewTopic(title, owner string, ownerId int, chatId int64) *Topic {
	topic := &Topic{
		Title:   title,
		Owner:   owner,
		OwnerId: ownerId,
		ChatId:  chatId,
	}

	common.DB.Create(topic)
	NewTopicReply(owner, ownerId, topic)
	return topic
}

func (topic *Topic) UpdateReplayContent(msgId int, content string) {
	common.DB.Model(&topic).Updates(map[string]interface{}{"replay_content": content, "last_message_id": msgId})
}

// 获取最大楼层数量
func (topic *Topic) MaxFloor() int {
	Lock.Lock()
	defer Lock.Unlock()
	var count = 0
	common.DB.Model(&TopicReply{}).Where("topic_id = ?", topic.ID).Count(&count)

	count++

	return count
}

func FindTopicByIdInChat(chatId int64, topicId uint) *Topic {
	Lock2.Lock()
	defer Lock2.Unlock()
	topic := &Topic{}
	common.DB.Where("ID = ? and chat_id = ?", topicId, chatId).First(&topic)

	if topic.ID == 0 {
		return nil
	}

	return topic
}
