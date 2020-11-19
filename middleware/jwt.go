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

func handleErr(c *gin.Context) {
	if utils.ReqFromHTML(c) {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	c.AbortWithStatusJSON(http.StatusForbidden, utils.JSON(utils.AuthenError, "", nil))
	return
}

// Auth 用户认证中间件
var Auth gin.HandlerFunc = func(c *gin.Context) {
	log.Printf("auth middleware")
	// token, err := getToken(c)
	// if err != nil || token == "" {
	// 	handleErr(c)
	// 	return
	// }
	// user, err := parseToken(token)
	// if err != nil {
	// 	handleErr(c)
	// 	return
	// }
	// fmt.Print(user)
	c.Next()
}

const (
	// SECRET SECRET
	SECRET string = "asdhjsdhhdhdhdhsasd"
)

func getToken(c *gin.Context) (token string, err error) {
	token = c.Query("token")
	if token != "" {
		return
	}
	token = c.PostForm("token")
	if token != "" {
		return
	}

	token = c.Request.FormValue("token")
	if token != "" {
		return
	}
	// var _token struct {
	// 	token string
	// }
	// if err = c.Copy().BindJSON(&_token); err != nil {
	// 	token = _token.token
	// 	if token != "" {
	// 		return
	// 	}
	// }
	token = c.Request.Header.Get("token")
	return
}

func createToken(u models.UserModel) (tokenstr string, err error) {
	//自定义claim
	claim := jwt.MapClaims{
		"id":   u.ID,
		"nick": u.Nick,
		"nbf":  time.Now().Unix(),            //指定时间之前 token不可用
		"iat":  time.Now().Unix(),            //签发时间
		"exp":  time.Now().Unix() + 60*60*24, //过期时间 24小时
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenstr, err = token.SignedString(SECRET)
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
func parseToken(tokenss string) (user *models.UserModel, err error) {
	token, err := jwt.Parse(tokenss, secret())
	if err != nil {
		err = errors.New("解析token出错")
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}

	user.ID = claim["id"].(uint) // uint64(claim["id"].(float64))
	user.Nick = claim["nick"].(string)

	exp := int64(claim["exp"].(float64))
	fmt.Println(user, "过期时间=====", time.Unix(exp, 0).Format("2001-01-02 15:04:05"))
	return
}
