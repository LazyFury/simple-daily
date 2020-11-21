package routes

import (
	"net/http"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProjectLog 项目日志
type ProjectLog struct{}

// Detail 日志详情
func (p *ProjectLog) Detail(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		panic("id不可为空")
	}

	log := &models.ProjectLogModel{}
	err := models.GetObjectOrNotFound(log, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		panic(utils.JSON(utils.NotFound, "", err))
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("", log))
}

// AddPage 添加日志
func (p *ProjectLog) AddPage(c *gin.Context) {
	projectID, _ := c.Params.Get("id")
	c.HTML(http.StatusOK, "project/log/add.tmpl", map[string]interface{}{
		"projectID": projectID,
	})
}

// UpdatePage 更新日志页面
func (p *ProjectLog) UpdatePage(c *gin.Context) {
	lid, _ := c.Params.Get("lid")

	if lid == "" {
		panic("请输入日志id")
	}

	log := &models.ProjectLogModel{}
	if err := models.GetObjectOrNotFound(log, map[string]interface{}{
		"id": lid,
	}); err != nil {
		panic(err)
	}
	c.HTML(http.StatusOK, "project/log/update.tmpl", map[string]interface{}{
		"log": log,
	})
}

//Add 添加项目日志
func (p *ProjectLog) Add(c *gin.Context) {
	log := &models.ProjectLogModel{}
	if err := c.ShouldBindJSON(log); err != nil {
		panic(utils.JSONError("绑定参数错误", err))
	}

	if err := log.Validator(); err != nil {
		panic(err)
	}

	// gorm 事务操作
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 保存日志
		if err := tx.Create(log).Error; err != nil {
			return err
		}

		// 查找项目
		project := &models.ProjectModel{Model: models.Model{ID: log.ProjectID}}
		if err := tx.Where(project).First(project).Error; err != nil {
			return err
		}

		// 更新项目
		project.Progress += log.PlusProgress
		if err := tx.Save(project).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("保存成功", log))
}

// Update 更新
func (p *ProjectLog) Update(c *gin.Context) {
	lid, _ := c.Params.Get("lid")

	if lid == "" {
		panic("日志id不可空")
	}
	log := &models.ProjectLogModel{}
	project := &models.ProjectModel{}
	// 查找日志
	if err := models.GetObjectOrNotFound(log, map[string]interface{}{
		"id": lid,
	}); err != nil {
		panic(err)
	}

	// 找到项目
	if err := models.DB.Where(map[string]interface{}{
		"id": log.ProjectID,
	}).First(project).Error; err != nil {
		panic(err)
	}
	// 减去旧的进度
	project.Progress -= log.PlusProgress

	// 绑定更新
	if err := c.ShouldBind(log); err != nil {
		panic(err)
	}
	// 验证更新
	if err := log.Validator(); err != nil {
		panic(err)
	}
	// 更新日志
	if err := models.DB.Save(log).Error; err != nil {
		panic(err)
	}
	// 增加进度 更新项目
	project.Progress += log.PlusProgress
	if err := models.DB.Save(project).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("更新成功", nil))
}
