package routes

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/Treblex/simple-daily/utils/sha"
	"github.com/gin-gonic/gin"
)

// User User
type User struct{}

// Index 用户列表
func (u *User) Index(c *gin.Context) {
	defer utils.GinRecover(c)
	page, size := models.GetPagingParams(c)
	users := &[]models.UserModel{}
	usersModel := models.GetObjectsOrEmpty(users, nil)
	if err := usersModel.Paging(page, size); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, utils.JSONSuccess("", usersModel.Result))
}

// Login 登录
func (u *User) Login(c *gin.Context) {
	defer utils.GinRecover(c)
	nick, ok := c.GetPostForm("nick")
	if !ok {
		panic("请输入用户名")
	}
	password, ok := c.GetPostForm("password")
	if !ok {
		panic("请输入密码")
	}

	user := &models.UserModel{}
	if err := models.DB.Where(map[string]interface{}{
		"nick": nick,
	}).Or(map[string]interface{}{
		"email": nick,
	}).First(user).Error; err != nil {
		panic("用户不存在")
	}
	hex := md5.Sum([]byte(sha.EnCode(password)))
	if user.Password == fmt.Sprintf("%x", hex) {
		c.JSON(http.StatusOK, utils.JSONSuccess("登陆成功", nil))
		return
	}

	c.JSON(http.StatusForbidden, utils.JSONError("密码错误", nil))
}

// Add 注册用户
func (u *User) Add(c *gin.Context) {
	defer utils.GinRecover(c)

	user := &models.UserModel{}
	if err := c.ShouldBind(user); err != nil {
		panic(err)
	}

	if err := user.Validator(); err != nil {
		panic(err)
	}
	hex := md5.Sum([]byte(sha.EnCode(user.Password)))
	user.Password = fmt.Sprintf("%x", hex)

	if err := models.DB.Create(user).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("注册成功", user))
}

// Update 更新
func (u *User) Update(c *gin.Context) {
	defer utils.GinRecover(c)

	user := &models.UserModel{}
	if err := models.DB.Where(map[string]interface{}{
		"id": 0,
	}).First(user).Error; err != nil {
		panic("用户不存在")
	}

	if err := c.ShouldBind(user); err != nil {
		panic(err)
	}

	if err := user.Validator(); err != nil {
		panic(err)
	}

	if err := models.DB.Save(user).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, utils.JSONSuccess("更新成功", user))
}
