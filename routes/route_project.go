package routes

import (
	"fmt"
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
	page, size := models.GetPagingParams(c)
	user := c.MustGet("user").(*models.UserModel)
	type (
		projectsType struct {
			*models.ProjectModel
			Favorited bool
		}
	)
	projects := &[]projectsType{}

	favorite := &models.FavoriteProjectModel{} //收藏的项目
	project := &models.ProjectModel{}          //项目
	projectModel := models.GetObjectsOrEmpty(projects, nil, func(db *gorm.DB) *gorm.DB {
		return db.Table(project.TableName()).Where(map[string]interface{}{
			"user_id": user.ID,
		}).Order("favorited desc,updated_at desc,id desc").Joins(fmt.Sprintf(
			"left join (select true favorited,project_id,user_id f_user_id from `%s` where `deleted_at` IS NULL) f on f.`project_id`=`%s`.`id` and f.`f_user_id`=`%s`.`user_id`",
			favorite.TableName(),
			project.TableName(), project.TableName(),
		))
	})

	if err := projectModel.Paging(page, size, func(db *gorm.DB) *gorm.DB {
		return db.Select([]string{"*"})
	}); err != nil {
		panic(err)
	}
	c.HTML(http.StatusOK, "project/index.tmpl", map[string]interface{}{
		"projects":   projects,
		"pagination": projectModel.Pagination,
		"user":       c.MustGet("user").(*models.UserModel),
	})
}

// Detail 项目详情
func (p *Project) Detail(c *gin.Context) {
	user := c.MustGet("user").(*models.UserModel)

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
	if err := models.GetObjectOrNotFound(project, map[string]interface{}{
		"id":      id,
		"user_id": user.ID,
	}); err != nil {
		panic(utils.JSON(utils.NotFound, "", err))
	}

	// 查询日志
	logs := &[]models.ProjectLogModel{}

	logsModel := models.GetObjectsOrEmpty(logs, nil, func(db *gorm.DB) *gorm.DB {
		return db.Where(map[string]interface{}{
			"project_id": project.ID,
		}).Order("created_at desc,id desc")
	})

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

	var jobs []string //聚合工作内容
	plusProgress := 0 //增加的进度
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
	id, _ := c.Params.Get("id")
	if id == "" {
		panic("请输入项目id")
	}

	user := c.MustGet("user").(*models.UserModel)

	project := &models.ProjectModel{}
	if err := models.GetObjectOrNotFound(project, map[string]interface{}{
		"id":      id,
		"user_id": user.ID,
	}); err != nil {
		panic(utils.JSON(utils.NotFound, "", nil))
	}

	c.HTML(http.StatusOK, "project/update.tmpl", map[string]interface{}{
		"project": project,
		"csrf":    c.MustGet("csrf").(string),
	})
}

// Add 添加项目
func (p *Project) Add(c *gin.Context) {
	user := c.MustGet("user").(*models.UserModel)
	project := &models.ProjectModel{UserID: user.ID}

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
	id, _ := c.Params.Get("id")
	if id == "" {
		panic("请输入项目id")
	}

	user := c.MustGet("user").(*models.UserModel)

	project := &models.ProjectModel{}

	if err := models.GetObjectOrNotFound(project, map[string]interface{}{
		"id":      id,
		"user_id": user.ID,
	}); err != nil {
		panic(utils.JSON(utils.NotFound, "", err))
	}

	if err := c.ShouldBind(project); err != nil {
		panic(utils.JSONError("绑定参数失败", err.Error()))
	}

	if err := project.Validator(); err != nil {
		panic(err)
	}

	project.UserID = user.ID

	row := models.DB.Save(project)
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
