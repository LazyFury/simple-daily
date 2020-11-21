package main

import (
	"html/template"
	"net/http"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/routes"
	"github.com/Treblex/simple-daily/tools"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	store := cookie.NewStore([]byte("scretsdd"))
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
	if err := models.Connect(`root:sukeaiya@tcp(localhost:1232)/daily?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai`); err != nil {
		panic(err)
	}

	// html模版
	_template := template.Must(tools.ParseGlob(template.New("base").Funcs(tools.TemplateFuncs), "templates", "*.tmpl"))
	g.SetHTMLTemplate(_template)

	// 注册路由
	routes.Start(g.Group(""))

	err := g.Run()
	if err != nil {
		panic(err)
	}
}
