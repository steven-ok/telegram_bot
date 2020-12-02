package topic

import (
	"github.com/jinzhu/gorm"
	"tg_bots/common"
)

type TopicReply struct {
	gorm.Model
	TopicId uint   `json:"topic_id" gorm:"desc:话题ID"`
	Owner   string `json:"owner" gorm:"desc:回复者用户名"`
	OwnerId int    `json:"owner_id" gorm:"desc:回复者ID"`
	Floor   int    `json:"floor" gorm:"desc:楼层"`
}

func NewTopicReply(owner string, ownerId int, topic *Topic) *TopicReply {
	// 查找我是否已经回复，如果已经回复，则我不可回复
	//
	replay := &TopicReply{}
	common.DB.Where("topic_id = ? and owner_id = ?", topic.ID, ownerId).First(&replay)

	if replay.ID != 0 {
		return nil
	}

	replay = &TopicReply{
		TopicId: topic.ID,
		Owner:   owner,
		OwnerId: ownerId,
		Floor:   topic.MaxFloor(),
	}

	common.DB.Create(replay)
	return replay
}
