package topic

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	topicModel "tg_bots/model/topic"
	"time"
)

type service struct {
	botApi *tgbotapi.BotAPI
}

func NewTopicService(botApi *tgbotapi.BotAPI) *service {
	return &service{botApi: botApi}
}

// 发起一个话题
func (ser *service) CreateTopic(command string, message *tgbotapi.Message, parameters ...string) {
	if len(parameters) >= 2 {
		title := safeStr(strings.Join(parameters[1:len(parameters)], " "))
		// 创建一个话题
		topic := topicModel.NewTopic(title, message.From.String(), message.From.ID, message.Chat.ID)

		//title := message.Text)
		title = fmt.Sprintf("[@`%s`](tg://user?id=%d) 发起一个接龙: 回复ID  *\\#%d*\\ \n\n*%s*\n\n1\\.\t@%s", message.From.String(), message.From.ID, topic.ID, title, message.From.String())

		sendMsg := ser.sendMsg(message.Chat.ID, title)
		topic.UpdateReplayContent(sendMsg.MessageID, title)

		ser.botApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    message.Chat.ID,
			MessageID: message.MessageID,
		})

		ser.advocates(message.Chat.ID, topic.ID)
	}
}

// 回复一个话题
func (ser *service) ReplyToTopic(command string, message *tgbotapi.Message, parameters ...string) {
	// 创建一个话题
	topicId, err := strconv.Atoi(parameters[1])

	if err != nil {
		ser.notice(message, fmt.Sprintf("[@`%s`](tg://user?id=%d)请回复正确的话题ID, [%s]非法！", message.From.String(), message.From.ID, strings.Join(parameters[1:], "")))
		return
	}

	topic := topicModel.FindTopicByIdInChat(message.Chat.ID, uint(topicId))
	if topic != nil {
		replay := topicModel.NewTopicReply(message.From.String(), message.From.ID, topic)
		fmt.Println(message.From)
		if replay != nil {
			text := topic.ReplayContent + fmt.Sprintf("\n%d.\t@%s", replay.Floor, message.From.String())
			// 删除上一条记录
			ser.botApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    message.Chat.ID,
				MessageID: topic.LastMessageId,
			})

			topic.UpdateReplayContent(ser.sendMsg(message.Chat.ID, text).MessageID, text)
		} else {
			ser.notice(message, fmt.Sprintf("[@`%s`](tg://user?id=%d)你已经回复过此话题 \\#%d", message.From.String(), message.From.ID, topicId))
		}
	} else {
		ser.notice(message, fmt.Sprintf("[@`%s`](tg://user?id=%d)话题 **\\#%d** 不存在 ", message.From.String(), message.From.ID, topicId))
	}
}

// 随机出现拥护者
func (ser *service) advocates(chatID int64, topicID uint) {
	// 随机出现拥护者
	go func() {
		list := topicModel.TopicAdvocatesByChatId(chatID)
		for _, advocate := range list {
			select {
			case <-time.After(time.Duration(rand.Intn(5)) * time.Second):
				topic := topicModel.FindTopicByIdInChat(advocate.ChatId, topicID)
				if topic != nil {
					replay := topicModel.NewTopicReply(advocate.Owner, advocate.OwnerId, topic)
					if replay != nil {
						text := topic.ReplayContent + fmt.Sprintf("\n%d\\.\t[@`%s`](tg://user?id=%d)", replay.Floor, advocate.Owner, advocate.OwnerId)
						// 删除上一条记录
						ser.botApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
							ChatID:    advocate.ChatId,
							MessageID: topic.LastMessageId,
						})
						sendMsg := ser.sendMsg(chatID, text)
						topic.UpdateReplayContent(sendMsg.MessageID, text)
					}
				}
			}
		}
	}()
}

func safeStr(text string) string {
	safeReg := regexp.MustCompile("([_\\[\\]()~>#+-=|{}`.!\\*])")
	return safeReg.ReplaceAllString(text, "\\$1")
}

func (ser *service) sendMsg(chatId int64, text string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "MarkdownV2"
	message, _ := ser.botApi.Send(msg)
	return message
}

func (ser *service) notice(message *tgbotapi.Message, text string) {
	go func() {
		sendMsg := ser.sendMsg(message.Chat.ID, text)
		<-time.After(7 * time.Second)
		ser.botApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    message.Chat.ID,
			MessageID: message.MessageID,
		})
		ser.botApi.DeleteMessage(tgbotapi.DeleteMessageConfig{
			ChatID:    message.Chat.ID,
			MessageID: sendMsg.MessageID,
		})
	}()
}
