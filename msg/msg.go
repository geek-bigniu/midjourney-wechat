package msg

import (
	"github.com/eatmoreapple/openwechat"
	"log"
	"mj-wechat-bot/api"
	"mj-wechat-bot/db"
	"strconv"
	"strings"
)

var (
	enableGroup []string
	Redis       *db.RedisUtil
)

func OnMessage(msg *openwechat.Message) {
	msgId := strconv.FormatInt(msg.NewMsgId, 10)
	log.Printf("msgId:%s", msgId)
	// 如果是文本消息, 并且内容为"ping", 则回复"pong"
	//if msg.IsText() && msg.Content == "ping" {
	//
	//	msg.ReplyText("pong")
	//}

	//if msg.IsPicture() {
	//	picture, err := msg.GetPicture()
	//	if err != nil {
	//		log.Printf("获取图片失败:%v", err)
	//		return
	//	}
	//	if picture != nil {
	//		location := picture.Request.URL.String()
	//		if err != nil {
	//			log.Printf("获取图片地址失败:%v", err)
	//			return
	//		}
	//		log.Printf("图片信息:%v", location)
	//		log.Printf("图片信息:%v", picture)
	//	}
	//
	//}

	if msg.IsSendByFriend() {
		// 获取发送用户信息
		sender, err := msg.Sender()
		if err == nil {
			log.Printf("%s", "==================收到信息====================")
			log.Printf("UserID: %s", sender.ID())
			log.Printf("NickName: %v", sender.NickName)
			log.Printf("MsgId: %v", msg.NewMsgId)
			log.Printf("Content: %v", msg.Content)
			log.Printf("%s", "==================信息结束====================\n\n")
		}
	}
	if msg.IsSendByGroup() {
		//群组信息
		sender, err := msg.Sender()
		if err == nil {
			//群组内发言的用户信息
			senderUser, err := msg.SenderInGroup()
			if err == nil {
				log.Printf("%s", "==================收到信息====================")
				log.Printf("GroupID: %s", sender.ID())
				log.Printf("GroupNickName: %v", senderUser.NickName)
				log.Printf("MsgId: %v", msg.NewMsgId)
				log.Printf("Content: %v", msg.Content)
				log.Printf("%s", "==================信息结束====================\n\n")
			}

		}

		//log.Printf("isOnwer: %v,NickName: %s,UserName: %s,ID :%s,Content: %s", sender.IsOwner, sender.NickName, sender.UserName, msg.Content)
	}

	//if sender.Uin == 0 {
	//	return
	//}
	log.Printf("消息类型:%v", msg.MsgType)
	if !msg.IsText() {
		log.Printf("非文本消息")
		return
	}
	if !msg.IsSendByFriend() && !msg.IsSendByGroup() {
		log.Printf("非好友和群组消息,忽略")
		return
	}
	realMsg := strings.TrimSpace(msg.Content)
	log.Println("收到消息:", realMsg)
	msg.AsRead()
	log.Println("收到消息1:", realMsg)
	log.Printf("commands:%v", Commands)
	for pre, command := range Commands {
		hasPrefix := strings.HasPrefix(realMsg, pre)
		log.Printf("判断命令:%s 结果:%v", pre, hasPrefix)
		if hasPrefix {
			log.Printf("开始设置NX:%s", msgId)
			if !api.CheckAPI(msgId) {
				log.Printf("消息已被处理，跳过")
				return
			}
			// 创建结构体实例
			impl := &Impl{
				msg:     msg,
				realMsg: realMsg,
			}
			impl.call(pre, command)
			return
		}
	}

	return
}
