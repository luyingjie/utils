package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	// "fmt"
	// "strings"
	// "github.com/dustin/go-humanize"
)

// DownloadFile : download file会将url下载到本地文件，它会在下载时写入，而不是将整个文件加载到内存中。
// 将数据流式传输到文件中，而不必将其全部加载到内存中, 因此大量小文件比较适合。
func DownloadFile(filepath string, url string) error {

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
	return err
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

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer file.Close()
		newFile, err := os.Create("./file/" + head.Filename)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer newFile.Close()

		_, err1 := io.Copy(newFile, file)
		if err1 != nil {
			io.WriteString(w, err1.Error())
			return
		}

		io.WriteString(w, "成功")

		http.Redirect(w, r, "/file/upload", http.StatusFound)
	}
}

func UploadPassHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		url := "http://192.168.181.4:8999/file/upload"
		file, head, err := r.FormFile("file")
		if err != nil {
			io.WriteString(w, "1-"+err.Error())
			return
		}
		defer file.Close()
		// 将接收到的 File 类型文件转成字节流
		// byte, err := ioutil.ReadAll(file)
		// 将磁盘文件转成字节流
		// byte, err := ioutil.ReadFile("./file/" + head.Filename)
		// if err != nil {
		// 	io.WriteString(w, err.Error())
		// }
		// -----------------------------------------方案 1---------------------------------------
		//创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
		// bodyBuf := &bytes.Buffer{}
		// bodyWriter := multipart.NewWriter(bodyBuf)

		// //创建第一个需要上传的文件,filepath.Base获取文件的名称
		// fileWriter, _ := bodyWriter.CreateFormFile("file", head.Filename)
		// //打开文件
		// // fd1, _ := os.Open(file1)
		// // defer fd1.Close()
		// //把第一个文件流写入到缓冲区里去
		// _, _ = io.Copy(fileWriter, file)

		// //这一句写入附加字段必须在_,_=io.Copy(fileWriter,fd)后面
		// // if len(param) != 0 {
		// // 	//param是一个一维的map结构
		// // 	for k, v := range param {
		// // 		bodyWriter.WriteField(k, v)
		// // 	}
		// // }
		// //获取请求Content-Type类型,后面有用
		// contentType := bodyWriter.FormDataContentType()
		// bodyWriter.Close()
		// //创建一个http客户端请求对象
		// client := &http.Client{}
		// //创建一个post请求
		// req, _ := http.NewRequest("POST", url, nil)
		// //设置请求头
		// // req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:47.0) Gecko/20100101 Firefox/47.0")
		// //这里的Content-Type值就是上面contentType的值
		// req.Header.Set("Content-Type", contentType)
		// //转换类型
		// req.Body = ioutil.NopCloser(bodyBuf)
		// //发送数据
		// data, _ := client.Do(req)
		// //读取请求返回的数据
		// bytes, _ := ioutil.ReadAll(data.Body)
		// defer data.Body.Close()
		// ----------------------------------------方案 2---------------------------------------------

		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		fileWriter, err := bodyWriter.CreateFormFile("file", head.Filename)
		if err != nil {
			io.WriteString(w, "2-"+err.Error())
		}

		_, err = io.Copy(fileWriter, file)
		if err != nil {
			io.WriteString(w, "3-"+err.Error())
		}

		contentType := bodyWriter.FormDataContentType()
		bodyWriter.Close()

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bodyBuf)
		if err != nil {
			io.WriteString(w, "4-"+err.Error())
		}
		// Add 和 Set都可以设置成功头信息
		req.Header.Add("content-type", contentType)
		// req.ContentLength = h.Size

		resp, err := client.Do(req)
		defer resp.Body.Close()

		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			io.WriteString(w, "5-"+err.Error())
		}

		// --------------------------------------------------------------------------------
		io.WriteString(w, "ok-"+string(resp_body))
		// http.Redirect(w, r, "/file/uploadpass", http.StatusFound)
	}
}

// func (c *UpFileController) Upload() {
// 	f, h, _ := c.GetFile("file")
// 	// path := "./" + h.Filename
// 	defer f.Close()
// 	// f.Close()

// 	c.SaveToFile("file", h.Filename)
// 	// h.content
// 	User := c.GetString("User")
// 	Key := c.GetString("Key")
// 	TaskId := c.GetString("TaskId")
// 	// fmt.Println(User)
// 	// fmt.Println(Key)
// 	url := operation.GetConf("FileServerUrl") + "upload/"
// 	b := httplib.Post(url)
// 	if TaskId != "" {
// 		b.Param("task_id", TaskId)
// 	}

// 	b.Header("Access-Key", Key)
// 	b.Header("User-Id", User)
// 	// b.Header("Content-Type", "multipart/form-data; boundary=asadadaddada")

// 	b.PostFile("file", h.Filename)
// 	// bt,err:=ioutil.ReadFile(h.Filename)
// 	// if err!=nil{
// 	//     fmt.Println("error 1")
// 	// }
// 	// bt,err:=ioutil.ReadAll(f)
// 	// if err!=nil{
// 	// }
// 	// defer func() {
// 	// 	os.Remove(h.Filename)
// 	// }()

// 	// b.Body(bt)
// 	var returnObj interface{}
// 	b.ToJSON(&returnObj)
// 	// str, err := b.String()
// 	// if err != nil {
// 	//     fmt.Println(err)
// 	// }

// 	c.Data["json"] = &returnObj

// 	c.ServeJSON()
// }

// // 文件得几个api这边服务没有做过多得控制，主要是前端保证正确性。
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

// func (c *GetFileListController) Get() {
// 	// http://127.0.0.1:8081/file/get_file_list?url=201801/1
// 	url := c.GetString("url")
// 	url = operation.GetConf("HistoryReport") + url
// 	files, _ := ioutil.ReadDir(url)
// 	fileList := make([]string, len(files))
// 	for i, file := range files {
// 		// if file.IsDir() {
// 		//     // listFile(myfolder + "/" + file.Name())
// 		//     fileList[i] = file.Name()
// 		// } else {
// 		//     // fmt.Println(myfolder + "/" + file.Name())
// 		// }
// 		fileList[i] = file.Name()
// 	}

// 	c.Data["json"] = fileList
// 	c.ServeJSON()
// }
