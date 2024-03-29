package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Treblex/simple-daily/models"
	"github.com/Treblex/simple-daily/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func handleErr(c *gin.Context, err string) {
	if utils.ReqFromHTML(c) {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, utils.JSON(utils.AuthedError, err, nil))
	return
}

// Auth 用户认证中间件
var Auth gin.HandlerFunc = func(c *gin.Context) {
	log.Printf("auth middleware")
	token, err := getToken(c)
	if err != nil || token == "" {
		handleErr(c, "")
		return
	}
	user, err := parseToken(token)
	if err != nil {
		handleErr(c, "解析token错误")
		return
	}
	c.Set("user", user)
}

const (
	// SECRET SECRET
	SECRET string = "asdhjsdhhdhdhdhsasd"
)

func getToken(c *gin.Context) (token string, err error) {
	// query
	token = c.Query("token")
	req := c.Request
	if token != "" {
		return
	}

	// post
	token = c.PostForm("token")
	if token != "" {
		return
	}

	token = req.FormValue("token")
	if token != "" {
		return
	}

	// header
	token = req.Header.Get("token")
	if token != "" {
		return
	}

	// cookie
	token, err = c.Cookie("token")
	if err != nil {
		return
	}

	// post json token内不做了，需要拷贝一份body，对性能有一些影响

	return
}

// CreateToken 创建token
func CreateToken(u models.UserModel) (token string, err error) {
	return CreateTokenMaxAge(u, int64(60*60*24))
}

// CreateTokenMaxAge 创建token
func CreateTokenMaxAge(u models.UserModel, maxAge int64) (tokens string, err error) {
	//自定义claim
	claim := jwt.MapClaims{
		"id":      u.ID,
		"nick":    u.Nick,
		"headPic": u.HeadPic,
		"nbf":     time.Now().Unix(),          //指定时间之前 token不可用
		"iat":     time.Now().Unix(),          //签发时间
		"exp":     time.Now().Unix() + maxAge, //过期时间 24小时
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err = token.SignedString([]byte(SECRET))
	return
}

// 解密token方法
func secret() jwt.Keyfunc {
	key := []byte(SECRET)
	return func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}
}

//ParseToken 解密token
func parseToken(tokens string) (user *models.UserModel, err error) {
	token, err := jwt.Parse(tokens, secret())
	if err != nil {
		err = errors.New("解析token出错")
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to map claim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	user = &models.UserModel{}
	user.ID = uint(claim["id"].(float64)) // uint64(claim["id"].(float64))
	user.Nick = claim["nick"].(string)
	user.HeadPic = claim["headPic"].(string)

	exp := int64(claim["exp"].(float64))
	fmt.Println(user.Nick, "过期时间=====", time.Unix(exp, 0).Format("2001-01-02 15:04:05"))
	return
}
