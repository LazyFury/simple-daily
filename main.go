package main

import (
	"html/template"
	"net/http"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/routes"
	"github.com/Treblex/simple-daily/tools"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	g.HandleMethodNotAllowed = true

	g.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, utils.JSON(http.StatusMethodNotAllowed, "方法不允许", nil))
	})

	g.NoRoute(func(c *gin.Context) {
		if utils.ReqFromHTML(c) {
			c.HTML(http.StatusMovedPermanently, "404.tmpl", nil)
			return
		}
		c.JSON(http.StatusNotFound, utils.JSON(http.StatusNotFound, "页面不存在", nil))
	})

	// 自定义验证器
	utils.RegValidator()

	// recover panic
	g.Use(gin.Recovery())

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
