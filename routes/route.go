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
	// 重置密码
	userRouter.GET("/reset-password", user.ResetPage)
	userRouter.POST("/reset-password", user.Reset)

	g.POST("/reg", user.Add)
	// 登录
	g.GET("/login", user.LoginPage)
	g.POST("/login", user.Login)
	// 登出
	g.GET("/logout", user.LogOut)
	// 忘记密码
	g.GET("/forgot", user.ForgotPage)
	g.POST("/forgot", user.Forgot)

}
