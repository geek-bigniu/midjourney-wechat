package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func GetImageUrlData(imageUrl string) (bool, io.Reader) {
	// 判断是否需要使用代理
	resp := &http.Response{}
	if config.Proxy.UseProxy {
		proxyUrl, err := url.Parse(config.Proxy.ProxyUrl)
		if err != nil {
			log.Fatalf("invalid proxy URL: %v", err)
		}

		transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		httpClient = &http.Client{Transport: transport, Timeout: 5 * time.Second}
		// 使用自定义的 HttpClient 发送请求

	} else {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := httpClient.Get(imageUrl)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to download image. StatusCode: %d\n", resp.StatusCode)
		return false, nil
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	imageReader := bytes.NewReader(imageData)
	return true, imageReader
}

func GetUserName(msg *openwechat.Message) (string, error) {
	if msg.IsSendByFriend() {
		// 获取发送用户信息
		sender, err := msg.Sender()
		if err != nil {
			return "", err
		}
		return sender.NickName, nil
	}
	if msg.IsSendByGroup() {
		//群组内发言的用户信息
		senderUser, err := msg.SenderInGroup()
		if err != nil {
			return "", err
		}
		return senderUser.NickName, nil
		//log.Printf("isOnwer: %v,NickName: %s,UserName: %s,ID :%s,Content: %s", sender.IsOwner, sender.NickName, sender.UserName, msg.Content)
	}
	return "", errors.New("为获取到")
}
