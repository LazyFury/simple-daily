package main

import (
	"html/template"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/routes"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Treblex/simple-daily/docs"
)

// @title simple-daily
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /
func main() {
	g := gin.Default()

	// 挂载静态文件
	g.Use(static.Serve("/static", static.LocalFile("static", false)))

	// 链接数据库
	if err := models.Connect(`root:sukeaiya@tcp(mysql:3306)/daily?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai`); err != nil {
		panic(err)
	}

	// 文档
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// html模版
	_template := template.Must(utils.ParseGlob(template.New("base").Funcs(utils.TemplateFuns), "templates", "*.tmpl"))
	g.SetHTMLTemplate(_template)

	// 注册路由
	routes.Start(g.Group(""))

	g.Run()
}
