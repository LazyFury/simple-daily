package models

import (
	"errors"
	"time"
)

// ProjectLogModel 项目日志
type ProjectLogModel struct {
	Model
	ProjectID    uint   `json:"project_id"`
	Content      string `json:"content"`
	PlusProgress int    `json:"plus_progress"`
}

var _ ModelType = &ProjectLogModel{}

// IsToday 在今天创建的
func (p *ProjectLogModel) IsToday() bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return p.CreatedAt.After(today)
}

// TableName 表名
func (p *ProjectLogModel) TableName() string {
	return "project_log_models"
}

// Validator 验证
func (p *ProjectLogModel) Validator() (err error) {

	if p.ProjectID == 0 {
		return errors.New("请输入项目id")
	}
	if p.Content == "" {
		return errors.New("请输入工作内容")
	}
	// if p.PlusProgress < 0 || p.PlusProgress > 100 {
	// 	return errors.New("进度应该在0-100之间")
	// }
	project := &ProjectModel{Model: Model{ID: p.ProjectID}}
	if DB.First(project).RowsAffected == 0 {
		return errors.New("项目不存在")
	}
	var progress = p.PlusProgress + project.Progress
	if progress > 100 || progress < 0 {
		return errors.New("项目进度应该在0-100之间")
	}

	return
}
