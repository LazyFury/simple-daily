package models

import (
	"errors"
	"strings"

	"github.com/Treblex/simple-daily/utils"
)

// ProjectModel 项目模型
type ProjectModel struct {
	Model
	Name               string            `json:"name" binding:"requiredParams" form:"name" gorm:"NOT NULL"`
	Describe           string            `json:"describe" form:"describe"`
	Start              utils.JSONTime    `json:"start"  form:"start"`
	ExpectEnd          utils.JSONTime    `json:"expect_end" form:"expect_end"`
	ActualDeliveryDate utils.JSONTime    `json:"actual_delivery_date" form:"actual_delivery_date"`
	Progress           int               `json:"progress" form:"progress"`
	Logs               []ProjectLogModel `json:"logs" gorm:"foreignKey:project_id"`
}

var _ ModelType = &ProjectModel{}

// TestVal 模板测试
func (p *ProjectModel) TestVal() string {
	return "hello"
}

// Validator 验证
func (p *ProjectModel) Validator() (err error) {
	p.Name = strings.Trim(p.Name, " ")
	if p.Name == "" {
		return errors.New("项目名不可空")
	}

	if p.Describe == "" {
		return errors.New("请输入项目描述")
	}

	if p.Start.IsZero() {
		return errors.New("请输入项目开始时间")
	}

	if p.ExpectEnd.IsZero() {
		return errors.New("请输入项目预计结束时间")
	}

	if p.ActualDeliveryDate.IsZero() {
		return errors.New("请输入项目交付时间")
	}

	if p.Progress < 0 || p.Progress > 100 {
		return errors.New("项目进度的值需要在0-100之间")
	}

	return
}
