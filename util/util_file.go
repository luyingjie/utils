package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	// "github.com/dustin/go-humanize"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func GetMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

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

//判断一个文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// DelFile : 删除文件
func DelFile(url string) error {
	err := os.Remove(url)

	if err != nil {
		return err
	}
	return nil
}

// SilenceDelFile : 静默删除文件
func SilenceDelFile(url string) {
	_ = os.Remove(url)
}

// GetFileSize : 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// GetFile : 获取文件
func GetFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 获取目录下的文件列表
func GetFileNameList(url string) []string {
	// http://127.0.0.1:8081/file/get_file_list?url=201801/1
	files, _ := ioutil.ReadDir(url)
	fileList := make([]string, len(files))
	for i, file := range files {
		// if file.IsDir() {
		//     // listFile(myfolder + "/" + file.Name())
		//     fileList[i] = file.Name()
		// } else {
		//     // fmt.Println(myfolder + "/" + file.Name())
		// }
		fileList[i] = file.Name()
	}

	return fileList
}

// GetBase64 : 将客户端的url转成base64
func GetBase64(url, _type string) (string, error) {
	str, err := GetFile(url)
	if err != nil {
		return "", err
	}
	base64 := "data:" + _type + ";base64," + base64.StdEncoding.EncodeToString(str)
	return base64, nil
}

// 保存文件
func SaveFile(url, name string, file multipart.File) error {
	newFile, err := os.Create(url + name)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		return err
	}
	return nil
}

// 将文件和表单发送出去, 返回服务器的返回body。 表单参数部分未测试。
func SendFile(url string, from map[string]string, name string, file multipart.File) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// bodyWriter.CreateFormField()
	fileWriter, err := bodyWriter.CreateFormFile("file", name)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", err
	}

	// 处理表单中的参数
	for k, v := range from {
		if err := bodyWriter.WriteField(k, v); err != nil {
			return "", err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return "", err
	}
	// Add 和 Set都可以设置成功头信息
	req.Header.Add("content-type", contentType)
	// req.ContentLength = h.Size

	resp, err := client.Do(req)
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(resp_body), nil
}

// DownloadFile : download file会将url下载到本地文件，它会在下载时写入，而不是将整个文件加载到内存中。
// 将数据流式传输到文件中，而不必将其全部加载到内存中, 因此大文件比较适合。
func DownloadFileToMem(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// WriteCounter : 计数器
// type WriteCounter struct {
// 	Total uint64
// }

// func (wc *WriteCounter) Write(p []byte) (int, error) {
// 	n := len(p)
// 	wc.Total += uint64(n)
// 	wc.PrintProgress()
// 	return n, nil
// }

// func (wc WriteCounter) PrintProgress() {
// 	fmt.Printf("\r%s", strings.Repeat(" ", 35))
// 	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
// }

// // DownloadFileCount : 可以传递计数器来跟踪进度。在下载时，我们还将文件另存为临时文件，因此在完全下载文件之前，我们不会覆盖有效文件。
// func DownloadFileCount(filepath string, url string) error {
// 	out, err := os.Create(filepath + ".tmp")
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		out.Close()
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	counter := &WriteCounter{}
// 	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
// 		out.Close()
// 		return err
// 	}
// 	fmt.Print("\n")
// 	out.Close()
// 	if err = os.Rename(filepath+".tmp", filepath); err != nil {
// 		return err
// 	}
// 	return nil
// }

// 获取文件，直接显示在浏览器中打开
// func (c *GetFileController) Get() {
// 	// http://127.0.0.1:8081/file/get_file?Year=2018&Month=01&BU=1&Name=test华创云平台资源服务明细账单-虚拟机维度.pdf
// 	year := c.GetString("year")
// 	month := c.GetString("month")
// 	bu := c.GetString("bu")
// 	name := c.GetString("name")

// 	url := operation.GetConf("HistoryReport") + year + month + "/" + bu + "/"

// 	pdfUrl := path.Join(url, name)
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
// 	c.Ctx.Output.Header("Content-Type", "application/pdf")
// 	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", pdfUrl))

// 	file, err := ioutil.ReadFile(pdfUrl)
// 	if err != nil {
// 		beego.Info("文件不存在")
// 		return
// 	}
// 	c.Ctx.WriteString(string(file))
// }

// 获取文件，直接获取下载的流，用于下载文件。
// func (c *DownFileController) Get() {
// 	// http://127.0.0.1:8081/file/get_file?Year=2018&Month=01&BU=1&Name=test华创云平台资源服务明细账单-虚拟机维度.pdf
// 	year := c.GetString("year")
// 	month := c.GetString("month")
// 	bu := c.GetString("bu")
// 	name := c.GetString("name")

// 	url := operation.GetConf("HistoryReport") + year + month + "/" + bu + "/" + name
// 	// b:=httplib.Get(url)
// 	// bt,err:=ioutil.ReadFile(url)
// 	// if err!=nil{
// 	//     // log.Fatal("read file err:",err)
// 	// }

// 	// b.ToFile("ddd.pdf")
// 	// fmt.Println("访问到了")
// 	c.Ctx.Output.Download(url, name)
// 	// c.Redirect("/static/img/logo.png",302)
// }
