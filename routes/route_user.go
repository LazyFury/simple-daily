package routes

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Treblex/simple-daily/config"
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
	page, size := models.GetPagingParams(c)
	users := &[]models.UserModel{}
	usersModel := models.GetObjectsOrEmpty(users, nil)
	if err := usersModel.Paging(page, size); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, utils.JSONSuccess("", usersModel.Result))
}

// ResetPage ResetPage
func (u *User) ResetPage(c *gin.Context) {
	csrf, _ := setCsrfKey(c)
	c.HTML(http.StatusOK, "login/reset.tmpl", map[string]interface{}{
		"csrf": csrf,
	})
}

// Reset 重设密码
func (u *User) Reset(c *gin.Context) {
	form := struct {
		Password        string `json:"password" form:"password" binding:"required"`
		PasswordConfirm string `json:"password_confirm" form:"password_confirm" binding:"required"`
	}{}
	user := c.MustGet("user").(*models.UserModel)
	if err := c.ShouldBind(&form); err != nil {
		panic(err)
	}

	if strings.Trim(form.Password, " ") == "" {
		panic("请输入密码")
	}

	if form.Password != form.PasswordConfirm {
		panic("两次输入的密码不相同")
	}

	if err := models.DB.Where(user).First(user).Error; err != nil {
		panic(err)
	}

	if user.Password == form.Password {
		panic("不可修改为和之前相同的密码")
	}
	user.Password = form.Password
	if err := user.Validator(); err != nil {
		panic(err)
	}

	password := sha.EnCode(form.Password)
	hex := md5.Sum([]byte(password))
	user.Password = fmt.Sprintf("%x", hex)

	if err := models.DB.Save(user).Error; err != nil {
		panic(err)
	}

	u.LogOut(c)
}

// ForgotPage 忘记密码
func (u *User) ForgotPage(c *gin.Context) {
	csrf, _ := setCsrfKey(c)
	c.HTML(http.StatusOK, "login/forgot.tmpl", map[string]interface{}{
		"csrf": csrf,
	})
}

// Forgot 忘记密码
func (u *User) Forgot(c *gin.Context) {
	form := struct {
		Email string `form:"email" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		panic(err)
	}

	if strings.Trim(form.Email, " ") == "" {
		panic("请输入用户邮箱")
	}
	// 查找用户
	user := &models.UserModel{Email: form.Email}
	if err := models.GetObjectOrNotFound(user, user); err != nil {
		panic("用户不存在")
	}

	// 创建临时token
	token, err := middleware.CreateTokenMaxAge(*user, int64(60*5)) //临时token 5分钟 重置密码成功之后删除重新登录
	if err != nil {
		panic(err)
	}

	// 发送重置密码邮件
	if err := config.Global.Mail.SendMail(
		"重置密码", []string{user.Email},
		fmt.Sprintf("重置密码链接: <a href='http://%s/users/reset-password?token=%s'>reset password</a>", c.Request.Host, token),
	); err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, utils.JSONSuccess("邮件发送成功！", nil))
}

// LogOut 用户登出
func (u *User) LogOut(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/login")
}

// LoginPage 登录页面
func (u *User) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login/login.tmpl", map[string]interface{}{
		"csrf": c.MustGet("csrf").(string),
	})
}

// Login 登录
func (u *User) Login(c *gin.Context) {
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

	c.JSON(http.StatusOK, utils.JSONError("密码错误", nil))
}

// Add 注册用户
func (u *User) Add(c *gin.Context) {

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

//UpdateProfile 更新用户信息
func (u *User) UpdateProfile(c *gin.Context) {
	user := c.MustGet("user").(*models.UserModel)
	if err := models.DB.First(user).Error; err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	c.HTML(http.StatusOK, "user/profile.tmpl", map[string]interface{}{
		"user": user,
		"csrf": c.MustGet("csrf").(string),
	})
}

// Update 更新
func (u *User) Update(c *gin.Context) {
	user := c.MustGet("user").(*models.UserModel)

	if err := models.DB.First(user).Error; err != nil {
		panic("用户不存在")
	}

	_, getFileErr := c.FormFile("file")
	if getFileErr == nil {
		path, err := Uploader.Default(c.Request)
		if err != nil {
			panic(err)
		}
		user.HeadPic = path
	}

	if err := c.ShouldBind(user); err != nil {
		panic(err)
	}

	log.Print(user.Password)

	if err := user.Validator(); err != nil {
		panic(err)
	}

	row := models.DB.Save(user)
	if err := row.Error; err != nil {
		panic(err)
	}
	token, err := middleware.CreateToken(*user)
	if err != nil {
		panic(err)
	}
	c.SetCookie("token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, utils.JSONSuccess("更新成功", nil))
}
