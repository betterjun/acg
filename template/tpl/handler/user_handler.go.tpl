package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"{{ .Name }}/tpl/errcode"
	"{{ .Name }}/tpl/params"
	"{{ .Name }}/tpl/svc"
)

// @Summary 获取服务器时间
// @Description 获取服务器时间
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} params.ServerTimeRsp
// @Router /v1/pub/servertime [GET]
// @ID GetServerTime
func GetServerTime(c *gin.Context) {
	now, ae := svc.GetServerTime()
	if ae != errcode.Success {
		c.JSON(http.StatusOK, errcode.Resp(ae, ae.Message()))
	} else {
		c.JSON(http.StatusOK, errcode.Resp(errcode.Success, params.ServerTimeRsp{ServerTime: now}))
	}
}
