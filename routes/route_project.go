package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// Project 项目
type Project struct{}

// Index 首页
func (p *Project) Index(c *gin.Context) {
	projects := &[]models.ProjectModel{}
	page, size := models.GetPagingParams(c)
	var midd models.Middleware = func(db *gorm.DB) *gorm.DB {
		return db.Order("updated_at desc,id desc")
	}
	projectModel := models.GetObjectsOrEmpty(projects, nil, midd)
	projectModel.Paging(page, size)
	c.HTML(http.StatusOK, "project/index.tmpl", map[string]interface{}{
		"projects":   projects,
		"pagination": projectModel.Pagination,
		"user":       c.MustGet("user").(*models.UserModel),
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
	logsModel := models.GetObjectsOrEmpty(
		logs,
		map[string]interface{}{
			"project_id": project.ID,
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		},
	)

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
			logsModel.Model = logsModel.Model.Where("`created_at` BETWEEN ? AND ?", startTime, endTime)
		}

	}

	if err := logsModel.All(); err != nil {
		panic(err)
	}

	jobs := []string{} //聚合工作内容
	plusProgress := 0  //增加的进度
	project.Logs = *logs

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
		"user":          c.MustGet("user").(*models.UserModel),
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
		panic(utils.JSONError(err.Error(), err))
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
func (p *Project) Update(c *gin.Context) {
	defer utils.GinRecover(c)
	id, _ := c.Params.Get("id")
	if id == "" {
		panic("请输入项目id")
	}

	db := models.DB
	project := &models.ProjectModel{}

	if err := db.Where(map[string]interface{}{
		"id": id,
	}).First(project).Error; err != nil {
		panic(utils.JSONError("项目不存在", err.Error()))
	}

	if err := c.ShouldBindJSON(project); err != nil {
		panic(utils.JSONError("绑定参数失败", err.Error()))
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
