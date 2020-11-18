package routes

import (
	"github.com/gin-gonic/gin"
)

// Start 注册路由
func Start(g *gin.RouterGroup) {
	// api := g.Group("/api")

	// 项目
	var p = Project{}
	var project = g.Group("/projects")

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
	logsRoute := g.Group("/project-logs")
	project.GET("/detail/:id/logs/add", pLogs.AddPage)
	project.POST("/detail/:id/logs/add", pLogs.Add)
	logsRoute.GET("/update/:lid", pLogs.UpdatePage)
	logsRoute.PUT("/update/:lid", pLogs.Update)
	logsRoute.GET("/detail/:id", pLogs.Detail)
}
