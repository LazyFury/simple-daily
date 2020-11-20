package upload

import (
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliOssConf AliOssConf
type AliOssConf struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"Access_key_secret"`
	URL             string `json:"url"`
}

// NewAliOssUploader 阿里云上传
func NewAliOssUploader(ali AliOssConf) *Uploader {
	return &Uploader{
		BaseDir:      "./static/oss",
		UploadMethod: ali.aliyunOssUpload,
		GetFile:      defaultGetFile,
	}
}

// AliyunOssUpload AliyunOssUpload
func (a *AliOssConf) aliyunOssUpload(name string, file io.Reader) (path string, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(fmt.Sprintf("Err:%x", err))
		}
	}()
	client, err := oss.New(a.Endpoint, a.AccessKeyID, a.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	bucket, err := client.Bucket("suke100")
	if err != nil {
		panic(err)
	}
	err = bucket.PutObject(name, file)
	if err != nil {
		panic(err)
	}
	path = fmt.Sprintf("%s%s", a.URL, name)
	return
}
