package routes

import (
	"net/http"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-gonic/gin"
)

// ProjectLog 项目日志
type ProjectLog struct{}

// AddPage 添加日志
func (p *ProjectLog) AddPage(c *gin.Context) {
	projectID, _ := c.Params.Get("id")
	c.HTML(http.StatusOK, "project/log/add.tmpl", map[string]interface{}{
		"projectID": projectID,
	})
}

// UpdatePage 更新日志页面
func (p *ProjectLog) UpdatePage(c *gin.Context) {
	defer utils.GinRecover(c)
	pid, _ := c.Params.Get("id")
	lid, _ := c.Params.Get("lid")
	if pid == "" {
		panic("请输入项目id")
	}
	if lid == "" {
		panic("请输入日志id")
	}

	log := &models.ProjectLogModel{}
	if err := models.DB.Where(map[string]interface{}{
		"id": lid,
	}).First(log).Error; err != nil {
		panic(err)
	}
	c.HTML(http.StatusOK, "project/log/update.tmpl", map[string]interface{}{
		"projectID": pid,
		"log":       log,
	})
}

//Add 添加项目日志
func (p *ProjectLog) Add(c *gin.Context) {
	defer utils.GinRecover(c)
	log := &models.ProjectLogModel{}
	if err := c.ShouldBindJSON(log); err != nil {
		panic(utils.JSONError("绑定参数错误", err))
	}

	if err := log.Validator(); err != nil {
		panic(err)
	}

	if err := models.DB.Create(log).Error; err != nil {
		panic(utils.JSONError("保存失败", err))
	}

	project := &models.ProjectModel{Model: models.Model{ID: log.ProjectID}}
	if models.DB.Find(project).RowsAffected == 0 {
		panic("没有找到对应的项目")
	}
	project.Progress += log.PlusProgress
	if err := models.DB.Save(project).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("保存成功", log))
}

// Update 更新
func (p *ProjectLog) Update(c *gin.Context) {
	defer utils.GinRecover(c)
	lid, _ := c.Params.Get("lid")
	pid, _ := c.Params.Get("id")
	if pid == "" {
		panic("请输入项目id")
	}
	if lid == "" {
		panic("日志id不可空")
	}
	log := &models.ProjectLogModel{}
	project := &models.ProjectModel{}
	// 查找
	if err := models.DB.Where(map[string]interface{}{
		"id":         lid,
		"project_id": pid,
	}).First(log).Error; err != nil {
		panic(err)
	}

	if err := models.DB.Where(map[string]interface{}{
		"id": pid,
	}).First(project).Error; err != nil {
		panic(err)
	}
	project.Progress -= log.PlusProgress

	// 绑定更新
	if err := c.ShouldBind(log); err != nil {
		panic(err)
	}

	if err := log.Validator(); err != nil {
		panic(err)
	}

	if err := models.DB.Save(log).Error; err != nil {
		panic(err)
	}

	project.Progress += log.PlusProgress
	if err := models.DB.Save(project).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("更新成功", nil))
}
