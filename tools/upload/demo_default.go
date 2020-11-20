package upload

// NewDefaultUploader 默认上传
func NewDefaultUploader() *Uploader {
	return &Uploader{
		BaseDir:      "./static/upload",
		UploadMethod: defaultUpload,
		GetFile:      defaultGetFile,
		MaxSize:      1024 * 1024 * 2,
	}
}
