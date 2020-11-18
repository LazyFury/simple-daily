package main

import (
	"html/template"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/routes"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()

	// 挂载静态文件
	g.Use(static.Serve("/static", static.LocalFile("static", false)))

	// 链接数据库
	if err := models.Connect(`root:sukeaiya@tcp(mysql:3306)/daily?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai`); err != nil {
		panic(err)
	}

	// html模版
	_template := template.Must(utils.ParseGlob(template.New("base").Funcs(utils.TemplateFuns), "templates", "*.tmpl"))
	g.SetHTMLTemplate(_template)

	// 注册路由
	routes.Start(g.Group(""))

	g.Run()
}
