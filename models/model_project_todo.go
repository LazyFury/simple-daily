package models

import (
	"gorm.io/gorm"
)

const (
	//ProjectTodoWait 等待处理
	ProjectTodoWait int = iota + 1
	// ProjectTodoWorking 工作中
	ProjectTodoWorking
	// ProjectTodoDone 完成
	ProjectTodoDone
)

var (
	projectTodoStatus = map[int]string{}
)

// ProjectTodoModel 待处理任务
type ProjectTodoModel struct {
	Model
	ProjectID uint           `json:"project_id" gorm:"comment:项目id"`
	UserID    uint           `json:"user_id" gorm:"comment:认领工作的用户"`
	Content   string         `json:"content"`
	Statuc    int            `json:"statuc" gorm:"comment:1-待认领,2-处理中,3-已完成"`
	DoneTime  gorm.DeletedAt `json:"done_time"`
	Done      bool           `json:"done" gorm:"-"`
}
