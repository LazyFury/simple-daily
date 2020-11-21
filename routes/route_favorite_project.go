package routes

import (
	"net/http"
	"strconv"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/gin-gonic/gin"
)

// FavoriteProject 收藏的项目
type FavoriteProject struct{}

// Add 收藏
func (f *FavoriteProject) Add(c *gin.Context) {
	user := c.MustGet("user").(*models.UserModel)
	projectID, _ := c.Params.Get("id")
	if projectID == "" {
		panic("项目id不可空")
	}
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		panic(err)
	}

	favorite := &models.FavoriteProjectModel{
		UserID:    user.ID,
		ProjectID: uint(pid),
	}
	if err := models.DB.Where(favorite).First(favorite).Error; err != nil {
		if err := models.DB.Create(favorite).Error; err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, utils.JSONSuccess("收藏成功", nil))
		return
	}

	if err := models.DB.Delete(favorite).Error; err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, utils.JSONSuccess("取消收藏", nil))
}
