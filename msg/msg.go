package msg

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
	"mj-wechat-bot/api"
	"mj-wechat-bot/db"
	"mj-wechat-bot/task"
	"strconv"
	"strings"
	"time"
)

var (
	enableGroup []string
	Redis       *db.RedisUtil
)

func OnMessage(msg *openwechat.Message) {
	msgId := strconv.FormatInt(msg.NewMsgId, 10)
	// 如果是文本消息, 并且内容为"ping", 则回复"pong"
	//if msg.IsText() && msg.Content == "ping" {
	//
	//	msg.ReplyText("pong")
	//}

	log.Printf("收到用户<%s>消息: %s,消息ID: %v", msg.FromUserName, msg.Content, msgId)

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
	//if msg.IsSendByGroup() {
	//	sender, _ := msg.SenderInGroup()
	//	sender1, _ := msg.Sender()
	//	//asGroup, _ := sender1.AsGroup()
	//	marshal, _ := json.Marshal(msg)
	//	log.Printf("Group: %s", marshal)
	//	log.Printf("GroupID: %s", sender1.ID())
	//	log.Printf("GroupNickName: %v", sender1.NickName)
	//	log.Printf("MsgId: %v", msg.NewMsgId)
	//	log.Printf("isOnwer: %v,NickName: %s,UserName: %s,ID :%s,Content: %s", sender.IsOwner, sender.NickName, sender.UserName, sender.ID(), msg.Content)
	//
	//}

	//if sender.Uin == 0 {
	//	return
	//}

	if msg.IsText() && (msg.IsSendByFriend() || msg.IsSendByGroup()) {
		//if msg.IsText() && (msg.IsSendByFriend()) {
		nowTime := time.Now().Unix()
		createTime := msg.CreateTime
		//log.Printf("收到消息时间:%v,当前时间:%v,时间差:%v", createTime, nowTime, nowTime-createTime)
		if nowTime-createTime > 10 {
			return
		}
		//https://wxapp.tc.qq.com/262/20304/stodownload?m=d4fc6f4d32185785a852940d6e5e7de2&filekey=30340201010420301e02020106040253480410d4fc6f4d32185785a852940d6e5e7de202022f27040d00000004627466730000000132&hy=SH&storeid=263203e1a0008fa23000000000000010600004f5053481f267b40b7966b65f&bizid=1023
		msg.AsRead()
		log.Printf("开始设置NX:%s", msgId)
		if !api.CheckAPI(msgId) {
			log.Printf("消息已被处理，跳过")
			return
		}
		//result, err := Redis.SetNX(msgId, "1", 10*time.Second)
		//if err != nil {
		//	log.Printf("错误信息:%v", err)
		//	log.Printf("消息已被处理，跳过")
		//	return
		//}
		//
		//if result {
		//	fmt.Println("Key set successfully")
		//} else {
		//	fmt.Println("Key already exists")
		//	log.Printf("消息已被处理，跳过1")
		//	return
		//}
		//// 获取 "@MJBOT " 部分的字节数
		//prefixBytes := []byte("@" + bot.CurrentNickName)
		//prefixSize := utf8.RuneCount(prefixBytes)
		//
		//// 获取字符串的前缀，并检查其是否为 "@MJBOT "
		//prefix := msg.Content[:prefixSize]
		//if prefix != "@"+bot.CurrentNickName {
		//	fmt.Println("Invalid prefix")
		//	msg.ReplyText("解析数据失败")
		//	return
		//}
		// 去除 "@MJBOT " 部分
		//realMsg := strings.TrimSpace(msg.Content[prefixSize:])
		//log.Println("收到消息:", realMsg)
		realMsg := strings.TrimSpace(msg.Content)
		//判断前缀是否为up命令
		if strings.HasPrefix(realMsg, "/help") {
			msg.ReplyText(fmt.Sprintf("%s\n当前支持的命令:\n%s\n\n%s\n%s\n%s\n\n%s\n%s",
				"------------------------",
				"/imagine 需要生成图片的描述",
				"/up      更新任务;示例: \"/up 任务id 参数\" \n参数可选:u1,u2,u3,u4,v1,v2,v3,v4",
				"   u1: u:选择第n张图片,n={1-4}",
				"   v1: v:生成与第n张图片相似的图片,n={1-4}",
				"/help    机器人帮助",
				"------------------------",
			))
			return
		}
		//判断前缀是否为up命令
		if strings.HasPrefix(realMsg, "/up") {
			realMsg = strings.ReplaceAll(realMsg, "/up", "")
			realMsg = strings.TrimSpace(realMsg)
			commands := strings.SplitN(realMsg, " ", 2)
			if len(commands) != 2 {
				msg.ReplyText("命令格式错误，示例:/up 任务id u1")
				return
			}
			taskId := strings.TrimSpace(commands[0])
			action := strings.ToLower(strings.TrimSpace(commands[1]))

			//判断action是否在指定字符串内
			switch action {
			case "u1", "u2", "u3", "u4", "v1", "v2", "v3", "v4":
				break
			default:
				msg.ReplyText("参数错误,可选参数:u1,u2,u3,u4,v1,v2,v3,v4")
				return
			}

			ok, newTaskId := api.TaskUpdate(taskId, action)
			if ok {
				msg.ReplyText("更新任务已经提交:" + newTaskId)
				log.Printf("任务已经提交:%s", newTaskId)
				task.AddTask(msg.FromUserName, newTaskId)
			} else {
				msg.ReplyText("任务创建失败")
			}
			return
		}
		//判断前缀是否为imagine命令
		if strings.HasPrefix(realMsg, "/imagine") {
			realMsg = strings.ReplaceAll(realMsg, "/imagine", "")
			realMsg = strings.TrimSpace(realMsg)
			ok, taskId := api.CreateMessage(realMsg)
			if ok {
				msg.ReplyText("任务已经提交:" + taskId)
				log.Printf("任务已经提交:%s", taskId)
				task.AddTask(msg.FromUserName, taskId)
			} else {
				msg.ReplyText("任务创建失败")
			}
			return
		}
	}
}
