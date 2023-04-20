package main

import (
	"fmt"
	"io"
	"log"
	"mj-wechat-bot/bot"
	"mj-wechat-bot/errorhandler"
	"mj-wechat-bot/msg"
	"os"
)

func main() {
	// 注册异常处理函数
	defer errorhandler.HandlePanic()
	// 创建日志文件
	file, err := os.OpenFile("wechat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("日志文件 wechat.log 创建失败: %v", err))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Close file error: %s", err)
		}
	}(file)
	// 创建一个多写器，同时将日志写入控制台和文件
	writers := []io.Writer{
		os.Stdout,
		file,
	}
	// 配置日志输出到文件
	log.SetOutput(io.MultiWriter(writers...))
	// 注册消息处理函数
	bot.Bot.MessageHandler = msg.OnMessage
	bot.StartBot()
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Bot.Block()
	// 等待操作系统信号或异常发生
	select {}

}
