package routes

import (
	"github.com/gin-gonic/gin"
)

// Start 注册路由
func Start(g *gin.RouterGroup) {
	// api := g.Group("/api")

	// 项目
	var p = Project{}
	var project = g.Group("/project")

	g.GET("/", p.Index)
	project.GET("/", p.Index)
	project.GET("/detail/:id", p.Detail)
	project.GET("/add", p.AddPage)
	project.POST("/add", p.Add)
	project.GET("/update/:id", p.UpdatePage)
	project.PUT("/update/:id", p.Update)
	project.DELETE("/del/:id", p.Delete)

	// 项目日志
	pLogs := &ProjectLog{}
	project.GET("/detail/:id/logs/add", pLogs.AddPage)
	project.POST("/detail/:id/logs/add", pLogs.Add)
	project.GET("/detail/:id/logs/update/:lid", pLogs.UpdatePage)
	project.PUT("/detail/:id/logs/update/:lid", pLogs.Update)
}
