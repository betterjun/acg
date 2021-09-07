package router

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
	"{{ .Name }}/pkg/cfg"
	"{{ .Name }}/pkg/logs"
	"{{ .Name }}/pkg/middleware"
	_ "{{ .Name }}/tpl/docs"
	"{{ .Name }}/tpl/handler"
)

func Init() error {
	r := gin.Default()

	// 配置了路径，才启用文件服务，部署时可以选择将前端的代码一起部署
	webPath := cfg.GetString("common.webpath")
	if len(webPath) > 0 {
		r.Static("/web", webPath)
	}

	// swagger文档
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 跨域
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowMethods:    []string{http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "X-Token"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials:       true,
		ExposeHeaders:          nil,
		MaxAge:                 0,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        true,
		AllowFiles:             true,
	}))

	v1 := r.Group("/v1")
	v1.Use(middleware.LogRequestMiddleware)
	v1.Use(middleware.LogResponseMiddleware)
	pub := v1.Group("/pub")
	{
		// 获取服务器时间
		pub.GET("/servertime", handler.GetServerTime)
	}

	if err := r.Run(fmt.Sprintf(":%s", cfg.GetString("common.port"))); err != nil {
		logs.Error("运行出错:%+v", zap.Error(err))
		return err
	}
	return nil
}
