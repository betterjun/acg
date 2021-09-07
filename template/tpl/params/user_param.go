package params

// 获取服务器时间响应
type ServerTimeRsp struct {
	// 服务器试卷
	ServerTime int64 `json:"server_time" form:"server_time"`
}
