package mysql

import (
	"fmt"
)

// Mysql Mysql
type Mysql struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Param    string `json:"param"`
}

// ToString ToString
func (m *Mysql) ToString() string {
	format := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True"
	return fmt.Sprintf(format, m.User, m.Password, m.Host, m.Port, m.Database) + m.Param
}
