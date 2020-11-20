package models

import (
	"errors"
	"regexp"
	"strings"
)

// UserModel 用户
type UserModel struct {
	Model
	Nick     string `json:"nick" form:"nick" gorm:"index;not null;unique"`
	HeadPic  string `json:"head_pic" form:"head_pic"`
	Email    string `json:"email" form:"email" gorm:"index;"`
	Password string `json:"password" gorm:"not null"`
}

var _ ModelType = &UserModel{}

// VerifyRepeatNickName 验证重复昵称
func (u *UserModel) VerifyRepeatNickName() error {
	user := &UserModel{Nick: u.Nick}
	if DB.Where(user).Not(map[string]interface{}{
		"id": u.ID,
	}).Find(user).RowsAffected > 0 {
		return errors.New("用户昵称已存在")
	}
	return nil
}

// VerifyRepeatEmail 验证重复邮箱
func (u *UserModel) VerifyRepeatEmail() error {
	if u.Email == "" {
		return nil
	}
	user := &UserModel{Email: u.Email}
	if DB.Where(user).Not(map[string]interface{}{
		"id": u.ID,
	}).First(user).RowsAffected > 0 {
		return errors.New("该邮箱已注册，是否需要找回密码？")
	}
	return nil
}

// Validator 验证
func (u *UserModel) Validator() error {
	u.Nick = strings.Trim(u.Nick, " ")
	if u.Nick == "" {
		return errors.New("请输入用户昵称")
	}
	if err := u.VerifyRepeatNickName(); err != nil {
		return err
	}
	if u.Password == "" {
		return errors.New("请设置您的密码")
	}
	tooWeak, err := regexp.MatchString(`^[\S]{6,48}$`, u.Password)
	if err != nil {
		return err
	}
	if !tooWeak {
		return errors.New("密码强度太低,6-24位字符")
	}

	if u.Email != "" {
		verify, err := regexp.MatchString(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, u.Email)
		if err != nil {
			return err
		}
		if !verify {
			return errors.New("邮箱格式不正确 email@xxx.com")
		}
		if err := u.VerifyRepeatEmail(); err != nil {
			return err
		}
	}

	return nil
}
