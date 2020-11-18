package models

// UserModel 用户
type UserModel struct {
	Model
	Nick     string `json:"nick"`
	HeadPic  string `json:"head_pic"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
