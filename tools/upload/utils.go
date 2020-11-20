package upload

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 在数组中
func inArray(arr []string, item string) (inArr bool) {
	index := -1
	item = strings.ToLower(item)
	for i, x := range arr {
		if item == x {
			index = i
		}
	}
	return index > -1
}

//GetDir 创建文件夹
func GetDir(path string, foderName string) (dir string, err error) {
	folder := filepath.Join(path, foderName)
	if _, err = os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return
		}
	}
	dir = folder
	return
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes 随机字符串
func RandStringBytes(n int) string {
	// 初始化随机数的资源库, 如果不执行这行, 不管运行多少次都返回同样的值
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
