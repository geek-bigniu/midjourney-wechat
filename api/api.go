package api

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"mj-wechat-bot/errorhandler"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	ApiUrl   string `yaml:"api_url"`
	Apikey   string `yaml:"api_key"`
	CheckUrl string `yaml:"check_url"`
}

var config API
var (
	createUrl     string
	taskUrl       string
	taskUpdateUrl string
)

func init() {
	// 注册异常处理函数
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
	createUrl = config.ApiUrl + "/imagine"
	taskUrl = config.ApiUrl + "/task"
	taskUpdateUrl = config.ApiUrl + "/task_update"

}

type Response struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

func CreateMessage(text string) (bool, string) {
	reqUrl, err := url.Parse(createUrl)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	query := reqUrl.Query()
	query.Set("prompt", text)
	reqUrl.RawQuery = query.Encode()
	body, err := DoGet(reqUrl)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	var response Response
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return false, ""
	}
	if response.Code != 1 {
		fmt.Println(response.Msg)
		return false, ""
	}
	return true, response.Data["task_id"].(string)
}

//查询任务状态
func QueryTaskStatus(taskID string) (bool, map[string]interface{}) {
	reqUrl, err := url.Parse(taskUrl)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	query := reqUrl.Query()
	query.Set("task_id", taskID)
	reqUrl.RawQuery = query.Encode()
	body, err := DoGet(reqUrl)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	log.Printf("任务【%s】返回结果 -> %s", taskID, body)
	var response Response
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return false, nil
	}
	if response.Code != 1 {
		fmt.Println(response.Msg)
		return false, nil
	}

	return true, response.Data
}

func TaskUpdate(taskId string, action string) (bool, string) {
	reqUrl, err := url.Parse(taskUpdateUrl)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	query := reqUrl.Query()
	query.Set("task_id", taskId)
	query.Set("action", action)
	reqUrl.RawQuery = query.Encode()
	log.Printf("reqUrl: %s", reqUrl.String())
	body, err := DoGet(reqUrl)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	var response Response
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		fmt.Println(err)
		return false, ""
	}
	if response.Code != 1 {
		fmt.Println(response.Msg)
		return false, ""
	}
	return true, response.Data["task_id"].(string)
}

func DoGet(reqUrl *url.URL) (string, error) {
	// 构建 HTTP GET 请求
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// 添加请求头
	req.Header.Add("apikey", config.Apikey)
	// 发送 HTTP GET 请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
}
