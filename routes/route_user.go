package routes

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/Treblex/simple-daily/middleware"
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

// LogOut 用户登出
func (u *User) LogOut(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/login")
}

// LoginPage 登录页面
func (u *User) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login/login.tmpl", nil)
}

// Login 登录
func (u *User) Login(c *gin.Context) {
	defer utils.GinRecover(c)
	postUser := struct {
		Nick     string `json:"nick"`
		Password string `json:"password"`
	}{}
	if err := c.Bind(&postUser); err != nil {
		panic(err)
	}
	if postUser.Nick == "" {
		panic("请输入用户名")
	}
	if postUser.Password == "" {
		panic("请输入密码")
	}

	user := &models.UserModel{}
	if err := models.DB.Where(map[string]interface{}{
		"nick": postUser.Nick,
	}).Or(map[string]interface{}{
		"email": postUser.Nick,
	}).First(user).Error; err != nil {
		panic("用户不存在")
	}

	hex := md5.Sum([]byte(sha.EnCode(postUser.Password)))
	if user.Password == fmt.Sprintf("%x", hex) {
		token, err := middleware.CreateToken(*user)
		if err != nil {
			panic(err)
		}
		c.SetCookie("token", token, 3600*24, "/", "", false, true)
		c.JSON(http.StatusOK, utils.JSONSuccess("登陆成功", token))
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
