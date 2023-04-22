package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mj-wechat-bot/errorhandler"
	"net/http"
	"net/smtp"
	"strings"
)

var (
	httpClient *http.Client
)

type ProxyConfig struct {
	UseProxy bool   `yaml:"use_proxy"`
	ProxyUrl string `yaml:"proxy_url"`
}

type Config struct {
	Smtp  Smtp        `yaml:"smtp"`
	Mail  MailConfig  `yaml:"mail"`
	Proxy ProxyConfig `yaml:"proxy"`
}

type Smtp struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MailConfig struct {
	Sender  string   `yaml:"sender"`
	To      []string `yaml:"to"`
	Subject string   `yaml:"subject"`
	Body    string   `yaml:"body"`
}

type Mail struct {
	SenderId string
	ToIds    []string
	Subject  string
	Body     string
}

var (
	config Config
	mail   Mail
)

func init() {
	defer errorhandler.HandlePanic()
	// Read configuration file.
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	// Unmarshal configuration.

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}

	// Send mail.
	mail = Mail{
		SenderId: config.Mail.Sender,
		ToIds:    config.Mail.To,
		Subject:  config.Mail.Subject,
		Body:     config.Mail.Body,
	}
}
func SendMail(url string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", config.Smtp.Username, config.Smtp.Password, config.Smtp.Host)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := strings.Join(mail.ToIds, ",")
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n<html><body>%s<br><img src=\"%s\"></body></html>", mail.SenderId, to, mail.Subject, mail.Body, url)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", config.Smtp.Host, config.Smtp.Port), auth, mail.SenderId, mail.ToIds, []byte(msg))
	if err != nil {
		log.Printf("Failed to send mail: %v", err)
	} else {
		log.Println("Mail sent successfully.")
	}
	return nil
}
