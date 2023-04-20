package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CheckAPI(id string) bool {
	// 创建一个 HTTP 客户端
	client := &http.Client{Timeout: 5 * time.Second}

	// 发起 GET 请求
	resp, err := client.Get(config.CheckUrl + "/check?id=" + id)
	if err != nil {
		return false
	}

	// 关闭响应体
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return false
	}
	result := &Result{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return false
	}

	// 判断 code 字段的值
	if result.Code == 1 {
		return true
	} else {
		return false
	}
}
