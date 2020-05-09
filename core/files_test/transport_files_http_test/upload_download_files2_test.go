package transport_files_http

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

/*
 示例：应用golang原生net/http模拟调用【服务端】上传和下载文件 add by syf 2020.5.8
     1,postman模拟客户端，调用/upload路由上传文件
     2,浏览器模拟客户端，从/files/filename路由下载文件

 golang自带"net/http"API参考:https://godoc.org/net/http#FileServer
 示例参考:
	 知乎：https://zhuanlan.zhihu.com/p/136410759
	 代码：https://github.com/zupzup/golang-http-file-upload-download/blob/master/main.go
*/

const maxUploadSize = 10 * 1024 * 1024 //设置最大上传文件大小为10MB
const uploadPath = "/tmp/shared"       //上传文件的目标路径
const filename = "rocksdb"             //文件名称，不带后缀

//测试二，模拟客户端向/upload路由发送文件，发送文件成功后服务即关闭
func TestHttpUploadFiles(t *testing.T) {
	//1,客户端(如postman)调用/upload上传文件
	http.HandleFunc("/upload", uploadFilesHandler())
	log.Print("Server started on localhost:8080, use /upload for uploading files")
	//设置服务监听端口，自定义
	srv := &http.Server{Addr: ":8080"}
	defer srv.Close()
	log.Fatal(srv.ListenAndServe())
}

//测试二，模拟客户端(如浏览器)从服务端路由/files/filename下载文件，下载完成，服务即关闭
func TestHttpDownloadFiles(t *testing.T) {
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))
	log.Print("Server started on localhost:8080, use /files/{fileName} for downloading")
	//服务端端口8081，自定义
	srv := &http.Server{Addr: ":8081"}
	defer srv.Close()
	log.Fatal(srv.ListenAndServe())
}

//文件上传
func uploadFilesHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}
		//parse and validate file and post parameters
		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		//获取上传的文件大小
		fileSize := fileHeader.Size
		fmt.Printf("File size (bytes): %v\n", fileSize)
		//校验文件上传的大小
		if fileSize > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		//读取文件字节数
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		//获取上传的文件类型,如application/zip
		detectedFileType := http.DetectContentType(fileBytes)
		fmt.Println("detectedFileType:", detectedFileType)
		//文件扩展名，如.zip。***注意***将main.go部署在centos上，会拿不到fileEndings的值。因此下方由filename+fileEndings[0]会出现异常。此时可直接指定扩展名。
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		fmt.Println("fileEndings:", fileEndings)
		//文件名+扩展名,***注意***部署在centos上，因缺少mime上述fileEndings拿到的是[]空数组。因此此处会报错。
		newPath := filepath.Join(uploadPath, filename+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)
		//检查目录是否存在，不存在则创建
		if !checkDirIsExists(uploadPath) {
			if err := os.MkdirAll(uploadPath, 0755); err != nil {
				renderError(w, "UPLOADPATH_NOT_EXISTS", http.StatusInternalServerError)
				return
			}
		}
		//写入文件
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		srv := startHttpServer()
		defer srv.Close()
		w.Write([]byte("SUCCESS"))
	})
}

//启动监听服务，供客户端下载文件
func startHttpServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	//returning reference so caller can call Shutdown()
	return srv
}

//返回错误的响应信息
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

//检查目录是否存在
func checkDirIsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
