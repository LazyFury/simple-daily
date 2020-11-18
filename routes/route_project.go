package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Project 项目
type Project struct{}

// Index 首页
func (p *Project) Index(c *gin.Context) {
	projects := &[]models.ProjectModel{}

	db := models.DB
	db.Order("updated_at desc").Find(projects)

	c.HTML(http.StatusOK, "project/index.tmpl", map[string]interface{}{
		"projects": projects,
	})
}

// Detail 项目详情
func (p *Project) Detail(c *gin.Context) {
	defer utils.GinRecover(c)

	id, _ := c.Params.Get("id")
	start, _ := c.GetQuery("start")
	end, _ := c.GetQuery("end")
	_type, _ := c.GetQuery("type")
	if id == "" {
		panic("请输入项目id")
	}

	// 初始化模型
	project := &models.ProjectModel{}
	// 查询详情
	if err := models.DB.Where(map[string]interface{}{
		"id": id,
	}).First(project).Error; err != nil {
		panic(err)
	}

	// 查询日志
	logs := &[]models.ProjectLogModel{}
	row := models.DB.Where(map[string]interface{}{
		"project_id": project.ID,
	}).Order("created_at desc")

	// 如果需要时间筛选
	if start != "" {
		endTime := time.Now()
		if end != "" {
			_endTime, err := time.Parse("2006-01-02", end)
			// 如果有结束时间 赋值，这里不能panic打断
			if err == nil {
				endTime = _endTime
			}

		}
		startTime, err := time.Parse("2006-01-02", start)
		if err == nil {
			row = row.Where("`created_at` BETWEEN ? AND ?", startTime, endTime)
		}

	}

	if err := row.Find(logs).Error; err != nil {
		panic(err)
	}

	project.Logs = *logs

	jobs := []string{} //聚合工作内容
	plusProgress := 0  //增加的进度

	for _, item := range *logs {
		jobs = append(jobs, item.Content)
		plusProgress += item.PlusProgress
	}

	// c.JSON(http.StatusOK, utils.JSONSuccess("", project))

	c.HTML(http.StatusOK, "project/detail.tmpl", map[string]interface{}{
		"project":       project,
		"jobs":          strings.Join(jobs, ","),
		"plus_progress": plusProgress,
		"isToday":       _type == "day",
		"isWeek":        _type == "week",
		"isMonth":       _type == "month",
	})

}

// AddPage 添加
func (p *Project) AddPage(c *gin.Context) {
	c.HTML(http.StatusOK, "project/add.tmpl", nil)
}

// UpdatePage 更新
func (p *Project) UpdatePage(c *gin.Context) {
	defer utils.GinRecover(c)
	id, _ := c.Params.Get("id")
	if id == "" {
		panic("请输入项目id")
	}

	project := &models.ProjectModel{}
	if err := models.DB.Where(map[string]interface{}{
		"id": id,
	}).First(project).Error; err != nil {
		panic("数据不存在")
	}

	c.HTML(http.StatusOK, "project/update.tmpl", map[string]interface{}{
		"project": project,
	})
}

// Add 添加项目
func (p *Project) Add(c *gin.Context) {
	defer utils.GinRecover(c)

	project := &models.ProjectModel{}

	if err := c.ShouldBindWith(project, binding.JSON); err != nil {
		panic(utils.JSONError("绑定参数失败", err))
	}

	if err := project.Validator(); err != nil {
		panic(utils.JSONError(err.Error(), nil))
	}

	db := models.DB
	if err := db.Create(project).Error; err != nil {
		panic(utils.JSONError("保存失败", err))
	}

	c.JSON(http.StatusOK, project)
}

// Update 更新项目
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param id path string true "项目ID"
// @Param body body models.ProjectModel true "更新"
// @Router /project/update/{id} [put]
func (p *Project) Update(c *gin.Context) {
	defer utils.GinRecover(c)
	id, _ := c.Params.Get("id")
	if id == "" {
		panic("请输入项目id")
	}

	db := models.DB
	project := &models.ProjectModel{}

	if db.Where(map[string]interface{}{
		"id": id,
	}).First(project).RowsAffected == 0 {
		panic("项目不存在")
	}

	if err := c.ShouldBindJSON(project); err != nil {
		panic(utils.JSONError("绑定参数失败", err))
	}

	if err := project.Validator(); err != nil {
		panic(err)
	}

	row := db.Save(project)
	if row.Error != nil {
		panic(row.Error)
	}

	if row.RowsAffected == 0 {
		panic("没有变动")
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("更新成功", nil))

}

// Delete 删除
func (p *Project) Delete(c *gin.Context) {
	defer utils.GinRecover(c)
	id, _ := c.Params.Get("id")
	if id == "" {
		panic("请传入项目id")
	}

	project := &models.ProjectModel{}
	if err := models.DB.Where(map[string]interface{}{
		"id": id,
	}).Delete(project).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("删除成功", nil))
}
