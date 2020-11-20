package upload

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// UploadMethod 拷贝文件到本地 或者 上传到oss的方法。返回url
type uploadMethod func(name string, file io.Reader) (url string, err error)

// 从request中获取到file
type getFile func(httpContext *http.Request) (file *multipart.FileHeader, err error)

// Uploader Uploader
type Uploader struct {
	BaseDir      string
	UploadMethod uploadMethod
	GetFile      getFile
}

// Default 默认上传类型和文件夹
func (u *Uploader) Default(httpContext *http.Request) (path string, err error) {
	return u.Custom(httpContext, []string{}, "default")
}

// Custom 自定义上传类型和目录
func (u *Uploader) Custom(httpContext *http.Request, acceptsExt []string, folder string) (url string, err error) {
	file, err := u.GetFile(httpContext)
	if err != nil {
		return
	}
	return u.uploadBase(file, acceptsExt, folder)
}

// OnlyAcceptsExt 限制类型 比如：仅图片
func (u *Uploader) OnlyAcceptsExt(httpContext *http.Request, acceptsExt []string, folder string) (url string, err error) {
	file, err := u.GetFile(httpContext)
	if err != nil {
		return
	}
	pathExt := path.Ext(file.Filename)
	// 自定义类型  覆盖前边的
	if inArray(acceptsExt, strings.Trim(pathExt, ".")) {
		return u.uploadBase(file, acceptsExt, folder)
	}
	err = errors.New("不允许上传这种类型的文件")
	return
}

// acceptsExt  这里是一个扩展到类型，默认到图片，视频 压缩包类型，已经写在默认方法中了
func (u *Uploader) uploadBase(file *multipart.FileHeader, acceptsExt []string, folderName string) (url string, err error) {
	pathExt := path.Ext(file.Filename)

	folder := ""
	// 如果符合类型，设定目录
	if inArray(AcceptsImgExt, strings.Trim(pathExt, ".")) {
		folder = "image"
	}
	if inArray(AcceptsVideoExt, strings.Trim(pathExt, ".")) {
		folder = "video"
	}
	if inArray(AcceptsAudioExt, strings.Trim(pathExt, ".")) {
		folder = "audio"
	}
	if inArray(AcceptsOtherFileExt, strings.Trim(pathExt, ".")) {
		folder = "file"
	}
	// 自定义类型  覆盖前边的
	if inArray(acceptsExt, strings.Trim(pathExt, ".")) {
		folder = folderName
	}
	// 如果不符合任何一种类型
	if folder == "" {
		err = errors.New("文件不合法")
		return
	}

	// 打开文件流
	src, err := file.Open()
	if err != nil {
		err = errors.New("打开文件失败")
		return

	}
	defer src.Close() //函数结束时自动关闭文件

	//创建文件夹
	dir, err := GetDir(path.Join(u.BaseDir, folder), time.Now().Format("2006_01_02"))
	if err != nil {
		err = errors.New("创建文件夹失败")
		return
	}

	// 随机文件名 + 文件后缀
	randName := RandStringBytes(32) + pathExt
	// Destination
	fileName := filepath.Join(dir, randName)

	url, err = u.UploadMethod(fileName, src)
	if url == "" {
		err = errors.New("上传失败")
		return
	}
	return
}
