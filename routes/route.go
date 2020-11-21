package routes

import (
	"github.com/Treblex/simple-daily/middleware"
	"github.com/Treblex/simple-daily/tools/upload"
	"github.com/gin-gonic/gin"
)

var (
	// Uploader 上传
	Uploader = upload.NewDefaultUploader()
)

// Start 注册路由
func Start(g *gin.RouterGroup) {
	// api := g.Group("/api")
	auth := g.Group("", middleware.Auth)
	// 项目
	var p = Project{}
	var project = auth.Group("/projects")

	auth.GET("/", p.Index)
	project.GET("/", p.Index)
	project.GET("/detail/:id", p.Detail)
	project.GET("/add", p.AddPage)
	project.POST("/add", p.Add)
	project.GET("/update/:id", p.UpdatePage)
	project.PUT("/update/:id", p.Update)
	project.DELETE("/del/:id", p.Delete)

	favorite := FavoriteProject{}
	project.POST("/favorite/:id", favorite.Add)

	// 项目日志
	pLogs := &ProjectLog{}
	logsRoute := auth.Group("/project-logs")
	project.GET("/detail/:id/logs/add", pLogs.AddPage)
	project.POST("/detail/:id/logs/add", pLogs.Add)
	logsRoute.GET("/update/:lid", pLogs.UpdatePage)
	logsRoute.PUT("/update/:lid", pLogs.Update)
	logsRoute.GET("/detail/:id", pLogs.Detail)

	// 用户
	user := &User{}
	userRouter := auth.Group("/users")
	userRouter.GET("/", user.Index)
	userRouter.GET("/profile", user.UpdateProfile)
	userRouter.POST("/profile", user.Update)

	g.POST("/reg", user.Add)
	g.GET("/login", user.LoginPage)
	g.POST("/login", user.Login)
	g.GET("/logout", user.LogOut)
}
