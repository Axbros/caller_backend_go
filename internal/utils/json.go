package utils

import (
	"encoding/json"
	"fmt"
)

// 封装的函数
func ConvertToJSON(event string, message string, data interface{}, key string) string {
	// 创建一个结构体来存储数据
	type Info struct {
		Event   string
		Message string
		Data    interface{}
		Key     string
	}
	info := Info{
		Event:   event,
		Message: message,
		Data:    data,
		Key:     key,
	}
	// 将结构体转换为 JSON 字符串
	jsonData, err := json.Marshal(info)
	if err != nil {
		fmt.Println("转换为 JSON 时出错:", err)
		return ""
	}
	return string(jsonData)
}

// 解析字符串为 JSON 的函数
func ParseTextToJSON(text string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
