package models

// FavoriteProjectModel 收藏的项目
type FavoriteProjectModel struct {
	Model
	UserID    uint `json:"user_id"`
	ProjectID uint `json:"project_id"`
}

var _ ModelType = &FavoriteProjectModel{}

// TableName 表名
func (f *FavoriteProjectModel) TableName() string {
	return "project_favorite_models"
}

// Validator Validator
func (f *FavoriteProjectModel) Validator() (err error) {
	return
}
