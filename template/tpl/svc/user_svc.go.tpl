package svc

import (
	"time"

	"{{ .Name }}/tpl/errcode"
)

// 获取服务器当前时间
func GetServerTime() (now int64, ae errcode.APIError) {
	return time.Now().Unix(), errcode.Success
}
