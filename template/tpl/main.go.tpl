package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"{{ .Name }}/pkg/cfg"
	"{{ .Name }}/pkg/db"
	"{{ .Name }}/pkg/logs"
	"{{ .Name }}/pkg/redis"
	"{{ .Name }}/tpl/model"
	"{{ .Name }}/tpl/router"
)

// @title 智谷星图考试系统 API
// @version 1.0
// @description pub开头的api不需要登录就可访问，pri开头的需要登录才能访问，访问pri开头的路由时，需要把登录返回的token放到header X-Token中，服务端会做鉴权。返回结果为json格式，包含字段code，data，msg三个字段 , 成功时，code为0，data有数据；失败时，code不为0，msg为错误消息。
// @termsOfService http://swagger.io/terms/

// @contact.name svinsight API Support
// @contact.url http://www.svinsight.com
// @contact.email svinsight@svinsight.com

// @license.name private
// @license.url http://www.svinsight.com/licenses/private.html

// @host localhost:16888
// @BasePath /
func main() {
	// 1 命令行参数解析
	mode := flag.String("m", "dev", "指定执行模式,只支持[dev|test|prod],默认是dev")
	flag.Parse()
	dev := true
	if *mode != "dev" {
		dev = false
	}
	if dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2 载入配置文件
	configFile := fmt.Sprintf("conf/%s.toml", *mode)
	fmt.Printf("使用的配置文件:%s..........\n", configFile)
	err := cfg.Initialize(configFile)
	if err != nil {
		fmt.Printf("读取配置文件失败[%s]: %s\n", configFile, err.Error())
		os.Exit(1)
	}

	// 3 初始化日志
	logs.Init(cfg.GetString("log.file"), cfg.GetString("log.level"), cfg.GetInt("log.maxSize"),
		cfg.GetInt("log.maxAge"), cfg.GetInt("log.maxBackup"), dev)

	// 4 sql数据库初始化
	err = InitDB()
	if err != nil {
		fmt.Printf("数据库初始化失败: %s\n", err.Error())
		logs.Info("数据库初始化失败", zap.Error(err))
		os.Exit(2)
	}
	model.Setup()

	// 5 nosql数据库初始化
	err = InitRedis()
	if err != nil {
		fmt.Printf("redis初始化失败: %s\n", err.Error())
		logs.Info("redis初始化失败", zap.Error(err))
		os.Exit(3)
	}

	// 6 初始化路由
	err = router.Init()
	if err != nil {
		fmt.Printf("web服务器启动失败: %s\n", err.Error())
		logs.Info("web服务器启动失败", zap.Error(err))
		os.Exit(4)
	}

	// 7 启动webserver

	logs.Info("程序已启动")

	// 阻塞
	select {}
}

// 连接数据库
func InitDB() (err error) {
	const SectionDB = "database"

	driver := cfg.GetString(SectionDB + ".driver")
	source := cfg.GetString(SectionDB + ".source")
	_, err = db.InitDefaultDB(source, driver, nil)
	return nil
}

// 连接nosql数据库
func InitRedis() error {
	const SectionDB = "redis"

	host := cfg.GetString(SectionDB + ".redis_host")
	auth := cfg.GetString(SectionDB + ".redis_auth")
	db := cfg.GetInt(SectionDB + ".redis_db")
	maxActive := cfg.GetInt(SectionDB + ".redis_max_active")
	maxIdle := cfg.GetInt(SectionDB + ".redis_max_idle")
	idleTimeout := cfg.GetInt(SectionDB + ".redis_idle_timeout")

	err := redis.Init(host, auth, db, maxActive, maxIdle, idleTimeout)
	if err != nil {
		return err
	}

	return nil
}
