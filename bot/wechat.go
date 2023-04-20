package bot

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"mj-wechat-bot/utils"
	"time"
)

var (
	Bot             *openwechat.Bot
	CurrentUser     *openwechat.Self
	CurrentNickName string
	reloadStorage   io.ReadWriteCloser
)

func init() {
	reloadStorage = openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	Bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
}
func Relogin(bot *openwechat.Bot) {
	//go bot.HotLogin(reloadStorage, &openwechat.RetryLoginOption{MaxRetryCount: 3})
}
func ConsoleQrCode(uuid string) {
	qrcodeUrl := "https://login.weixin.qq.com/l/" + uuid

	q, _ := qrcode.New(qrcodeUrl, qrcode.Low)
	fmt.Println(q.ToString(true))
	go utils.SendMail(qrcodeUrl)

}
func StartBot() {
	Bot.UUIDCallback = ConsoleQrCode
	Bot.LogoutCallBack = Relogin
	PushLogin()
	go checkLife()
	// 获取登陆的用户
	CurrentUser, err := Bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	CurrentNickName = CurrentUser.NickName
	// 获取所有的好友
	friends, err := CurrentUser.Friends()
	log.Printf("%v", friends)

	// 获取所有的群组
	groups, err := CurrentUser.Groups()
	log.Println(groups, err)
	Bot.Block()
}
func PrintlnQrcodeUrl(uuid string) {
	println("访问下面网址扫描二维码登录")
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	println(qrcodeUrl)

	err := utils.SendMail(qrcodeUrl)
	if err != nil {
		log.Printf("Failed to send mail: %v", err)
	} else {
		fmt.Println("Mail sent successfully.")
	}

}
func checkLife() {
	for {
		life := Bot.Alive()
		log.Printf("存活状态:%v", life)
		if !life {
			log.Printf("机器人已掉线，尝试重新登陆")
			PushLogin()

		}
		time.Sleep(3 * time.Second)
	}
}
func PushLogin() {
	time.Sleep(1 * time.Second)
	// 创建热存储容器对象

	Bot.SyncCheckCallback = func(resp openwechat.SyncCheckResponse) {

	}
	// 执行热登录
	err := Bot.HotLogin(reloadStorage)
	if err != nil {
		// 执行提示登录
		log.Printf("开始尝试直接登陆")
		err := Bot.PushLogin(reloadStorage, &openwechat.RetryLoginOption{MaxRetryCount: 3})
		if err != nil {
			log.Printf("登陆失败:%v", err)
			return
		}

	}
}
