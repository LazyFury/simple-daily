package models

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Model Model
type Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type middleware func(db *gorm.DB) *gorm.DB

// Objects List
type Objects struct {
	Obj   interface{}
	Model *gorm.DB
}

// All 全部数据
func (o *Objects) All() (err error) {
	o.Model.Find(o.Obj)
	return
}

// Paging 分页数据
func (o *Objects) Paging(page int, size int) (err error) {
	offset := size * (page - 1)
	row := o.Model.Offset(offset).Limit(size).Find(o.Obj)
	if row.Error != nil {
		err = row.Error
	}
	return
}

// GetObjectsOrEmpty 获取列表 \n
// 可选参数 middleware models.middleware 接收一个 *gorm.DB 返回 *gorm.DB
func GetObjectsOrEmpty(obj interface{}, query interface{}, args ...interface{}) *Objects {
	row := DB.Model(obj)
	// 可选参数
	for _, arg := range args {
		midd, ok := arg.(middleware)
		if ok {
			row = midd(row)
		}
	}
	return &Objects{
		Model: row,
		Obj:   obj,
	}
}

// GetParamsTryInt 字符串转数字
func GetParamsTryInt(val string, defaul int) int {
	num, err := strconv.Atoi(val)
	if err != nil {
		num = defaul
	}
	return num
}

// GetPagingParams 获取分页参数
func GetPagingParams(c *gin.Context) (page int, size int) {
	page = GetParamsTryInt(c.Query("page"), 1)
	size = GetParamsTryInt(c.Query("size"), 10)
	return
}
