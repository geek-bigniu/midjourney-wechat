package task

import (
	"bytes"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"golang.org/x/image/webp"
	"image/png"
	"io"
	"log"
	"mj-wechat-bot/api"
	"mj-wechat-bot/bot"
	"mj-wechat-bot/utils"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

var (
	taskIds = sync.Map{}
	wg      = sync.WaitGroup{}
)

type ImageMsg struct {
	taskId       string
	fromUserName string
	url          string
}

func init() {
	go Looper()
	go ImageSender()
}

var (
	count   = int64(0)
	msgChan = make(chan ImageMsg, 10)
	test    = sync.RWMutex{}
)

// AddTask 添加任务
func AddTask(fromUserName string, taskId string) {
	log.Printf("添加任务:%s", taskId)
	atomic.AddInt64(&count, 1)
	taskIds.Store(taskId, fromUserName)
}

func ImageSender() {
	for {
		select {
		case msg := <-msgChan:
			log.Printf("收到发送图片任务，开始发送图片")
			log.Printf("%v", msg)
			//判断任务是否已经成功发送并删除
			_, ok := taskIds.Load(msg.taskId)
			if !ok {
				continue
			}
			// 发送图片消息
			sendImage(msg.taskId, msg.fromUserName, msg.url)
			time.Sleep(3 * time.Second)
		}
	}
}

// Looper 任务循环
func Looper() {
	log.Printf("开始启动任务循环")
	for {
		log.Printf("任务数量:%d", count)
		taskIds.Range(func(taskId, _ any) bool {

			wg.Add(1)
			// 查询任务状态
			go QueryTaskStatus(taskId.(string))
			return true
		})
		wg.Wait()
		time.Sleep(10 * time.Second)
	}
}

// QueryTaskStatus 查询任务状态并发送图片消息
func QueryTaskStatus(taskId string) {
	// 查询任务状态
	ok, data := api.QueryTaskStatus(taskId)
	if ok {
		value, _ := taskIds.Load(taskId)
		fromUserName := value.(string)
		// 判断是否完成
		switch data["status"] {
		case "finished":
			url := data["image_url"].(string)

			msgChan <- ImageMsg{
				taskId:       taskId,
				fromUserName: fromUserName,
				url:          url,
			}

			break
		case "pending":
			// 任务未完成
			break
		case "wait":
			// 任务未完成
			break
		case "invalid params":
			// 任务参数错误
			msg := fmt.Sprintf("任务被拒绝，参数错误，请检查:%s,删除任务", taskId)
			failTask(taskId, fromUserName, msg)
			break
		case "banned":
			// 任务被封禁
			msg := fmt.Sprintf("任务被拒绝，可能包含违禁词:%s,删除任务", taskId)
			failTask(taskId, fromUserName, msg)
			break
		}

	}
	wg.Done()
}

func failTask(taskId string, fromUserName string, msg string) {
	req := bot.Bot.Storage.Request
	info := bot.Bot.Storage.LoginInfo
	log.Printf("req:%v,info:%v,bot.CurrentUser:%s,fromUserName:%s", req, info, bot.CurrentUser, fromUserName)
	// 获取登陆的用户
	CurrentUser := bot.CurrentUser
	_, err := bot.Bot.Caller.WebWxSendMsg(&openwechat.SendMessage{
		FromUserName: CurrentUser.UserName,
		ToUserName:   fromUserName,
		Content:      msg,
	}, info, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	//删除任务
	taskIds.Delete(taskId)
}

// 发送图片消息
func sendImage(taskId string, fromUserName string, url string) {
	// 发送图片消息
	ok, reader := utils.GetImageUrlData(url)
	// 通过 path.Ext 函数解析链接地址中的后缀名
	ext := path.Ext(url)
	// 根据后缀名判断是否是 webp 格式的图片
	if ext == ".webp" {
		image, err := webp.Decode(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 创建一个 PNG 格式的 io.Reader
		var pngReader io.Reader
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, image); err != nil {
			fmt.Printf("pngReader: %v", err)
			return
		}

		pngReader = bytes.NewReader(buf.Bytes())
		reader = pngReader

	}
	if !ok {
		log.Printf("发送图片失败(%s):%s", taskId, url)
		return
	}
	// 发送图片消息
	req := bot.Bot.Storage.Request
	info := bot.Bot.Storage.LoginInfo
	log.Printf("req:%v,info:%v,bot.CurrentUser:%s,fromUserName:%s", req, info, bot.CurrentUser, fromUserName)
	// 获取登陆的用户
	CurrentUser, err := bot.Bot.GetCurrentUser()
	if err != nil {
		fmt.Printf("获取当前登陆用户失败:%s", err)
		return
	}
	_, err = bot.Bot.Caller.WebWxSendImageMsg(reader, req, info, CurrentUser.UserName, fromUserName)
	if err != nil {
		fmt.Println(err)
		return
	}
	//完成任务
	log.Printf("发送图片完成:%s,删除任务", url)
	taskIds.Delete(taskId)
	atomic.AddInt64(&count, -1)
}
