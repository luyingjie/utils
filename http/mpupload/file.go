package mpupload

import (
	"mime/multipart"
)

type MultipartUploadModel struct {
	Name     string
	FileHead *multipart.FileHeader
	File     multipart.File
}

// func Mpupload(url string, param map[string]string, heads map[string]string, files []FileModel) string {

// }

// MultipartUpload : 分块上传文件
func MultipartUpload(FileHead *multipart.FileHeader, File multipart.File) {
	// 1. 获取文件
	// 2. 按文件唯一ID查看是否已经上传
	// 3. If 查询Redis是否有分块信息，
	//    No 就是进行分块， 并写入分块信息到Redis。
	//    Yes 就上传没有传完的部分。
}
