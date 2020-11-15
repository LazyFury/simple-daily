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
	g.Use(static.Serve("/static", static.LocalFile("static", false)))
	if err := models.Connect(); err != nil {
		panic(err)
	}

	_template := template.Must(utils.ParseGlob(template.New("base").Funcs(utils.TemplateFuns), "templates", "*.tmpl"))
	g.SetHTMLTemplate(_template)

	routes.Start(g.Group(""))

	g.Run()
}
