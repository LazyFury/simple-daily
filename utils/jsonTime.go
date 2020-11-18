package utils

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

var (
	// DefaultTimeLayout DefaultTimeLayout
	DefaultTimeLayout string = "2006年01月02日 15:04:05"
)

type (
	// JSONTime JSONTime
	JSONTime struct {
		time.Time
	}
)

// ToDate ToDate
func (t *JSONTime) ToDate() string {
	return t.Format("2006年01月02日")
}

// ToString ToString
func (t *JSONTime) ToString() string {
	return t.Format(t.localLayout())
}

func (t *JSONTime) localLayout() string {
	return "2006年01月02日 15:04:05"
}

// MarshalJSON json格式化时间的方法
func (t JSONTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t.Time).Format(t.localLayout()))
	// var stamp = fmt.Sprintf("%d", time.Time(t.Time).Unix())
	return []byte(stamp), nil
}

// UnmarshalJSON UnmarshalJSON
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	var str = string(b)
	fmt.Printf(str + "\n")
	tTime, err := time.Parse(`"`+t.localLayout()+`"`, str)
	if err == nil {
		t.Time = tTime
		return nil
	}

	log.Debugf("format default fail")

	tTime, err = time.Parse(`"2006-01-02"`, str)
	if err == nil {
		t.Time = tTime
		return nil
	}

	tTime, err = time.Parse(`"2006-01-02T15:04:05"`, str)
	if err == nil {
		t.Time = tTime
		return nil
	}
	return errors.New("时间格式错误")
}

// Value Value
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan Scan
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
