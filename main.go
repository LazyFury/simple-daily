package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Treblex/simple-daily/config"
	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/routes"
	"github.com/Treblex/simple-daily/tools"
	"github.com/Treblex/simple-daily/utils"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	sentryIniterr := sentry.Init(sentry.ClientOptions{
		Dsn: "https://af4ccc0783a64fbc8faa1d42b030a36d@o299150.ingest.sentry.io/5613819",
	})
	if sentryIniterr != nil {
		log.Fatalf("sendtry init err %v\n", sentryIniterr)
	}
	defer sentry.Flush(2 * time.Second)

	store := cookie.NewStore([]byte("secrets"))
	g.Use(sessions.Sessions("daily", store))

	g.HandleMethodNotAllowed = true

	g.NoMethod(func(c *gin.Context) {
		panic(utils.JSON(http.StatusMethodNotAllowed, "", nil))
	})

	g.NoRoute(func(c *gin.Context) {
		panic(utils.JSON(http.StatusNotFound, "", nil))
	})

	g.Use(gin.Logger())

	// recover panic
	g.Use(gin.Recovery())

	g.Use(func(c *gin.Context) {
		defer utils.GinRecover(c)
		c.Next()
	})

	// 自定义验证器
	utils.RegValidator()

	// 挂载静态文件
	g.Use(static.Serve("/static", static.LocalFile("static", false)))

	// 链接数据库
	if err := models.Connect(config.Global.Mysql.ToString()); err != nil {
		panic(err)
	}

	// html模版
	_template := template.Must(tools.ParseGlob(template.New("base").Funcs(tools.TemplateFuncs), "templates", "*.tmpl"))
	g.SetHTMLTemplate(_template)

	// 注册路由
	routes.Start(g.Group(""))

	// ico
	g.GET("/favicon.ico", func(c *gin.Context) {
		c.File("static/favicon.ico")
	})

	// 启动
	err := g.Run(fmt.Sprintf(":%d", config.Global.Port))
	if err != nil {
		panic(err)
	}
}
