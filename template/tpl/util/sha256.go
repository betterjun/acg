package util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"pkg/cfg"
	"pkg/logs"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func StringSha256(str string) string {
	sum := sha256.Sum256([]byte(str))
	hash := fmt.Sprintf("%x", sum)
	return hash
}

func HeaderCheck(c *gin.Context) error {
	router := c.FullPath()
	logs.Sugar().Info("router:", router)
	expire := c.GetHeader("expirets")
	logs.Sugar().Info("expire:", expire)
	ctis := c.GetHeader("ctis")
	logs.Sugar().Info("ctis:", ctis)
	scert := cfg.GetString("ctis.secret1")
	logs.Sugar().Info("scert:", scert)
	er := errors.New("非法请求")
	if expire == "" || ctis == "" {
		return er
	}
	sha := StringSha256(router + "_" + expire + "_" + scert)
	if sha != ctis {
		return er
	}
	expireTime, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		return er
	}
	if time.Now().Unix() >= expireTime {
		return er
	}
	return nil
}
