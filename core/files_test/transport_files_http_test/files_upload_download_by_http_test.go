package transport_files_http

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
   示例：在peer中启动该服务，用于接收外部golang请求，上传与下载文件 add by syf 2020.5.12
   项目需求： 拷备rocksdb目录
   逻辑设计： 1，/isReady将rocksdb目录打成.zip包，同时启动下载服务；并告知客户端已准备好。可以下载
			2，客户端调用http url下载rocksdb.zip文件
			3，客户端调用/upload上传文件，由客户端指定具体下载到哪个目录下
*/

var srv *http.Server

func TestFilesUploadAndDownload(t *testing.T) {
	http.HandleFunc("/isReady", isReady)             //准备工作
	http.HandleFunc("/upload", uploadFileHandler)    //上传文件
	http.HandleFunc("/shutdown", shutdownHttpServer) //关闭下载服务
	http.ListenAndServe(":8001", nil)                //设置服务监听端口,下述还有一个文件下载服务端口为8000
}

//准备工作
func isReady(w http.ResponseWriter, r *http.Request) {
	zipDir := "/tmp/zipserver" //打成.zip包所存储的路径及下载服务所指向的目录位置
	//zipFileName:="rocksdb0.zip"
	//zipFileName:="/tmp/peer0/rocks.db/rocksdb0.db"
	//此参数设置参见：http_send_form-data_test.go/TestHttpSendFormData
	filespath := r.FormValue("filespath")          //压缩后所保存的路径,可由客户端指定
	zipFileName := r.FormValue("dirname") + ".zip" //压缩的文件名称，可由客户端指定
	if !checkDirIsExist(zipDir) {
		if err := os.MkdirAll(zipDir, 0755); err != nil {
			renderErrors(w, "UPLOADPATH_NOT_EXISTS", http.StatusInternalServerError)
			return
		}
	}
	//压缩rocksdb【param1:压缩的.zip文件;param2:压缩后所保存的路径】
	if err := zipfiles(filepath.Join(zipDir, zipFileName), filespath); err != nil {
		w.Write([]byte("fail"))
	}
	//开启下载服务监听，客户端下载完成后，会发通知关闭下载服务
	go startingHttpServer(zipDir)
	w.Write([]byte("ok"))
}

//启动监听服务，供客户端下载文件
func startingHttpServer(monitorDir string) {
	srv = &http.Server{Addr: ":8000"}
	//defer srv.Close()
	fs := http.FileServer(http.Dir(monitorDir))
	http.Handle("/files/", http.StripPrefix("/files", fs))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	//time.Sleep(15*time.Minute) //15min后服务自动关闭
}

//关闭文件下载服务监听
func shutdownHttpServer(w http.ResponseWriter, r *http.Request) {
	srv.Close()
	w.Write([]byte("ok"))
}

//文件上传
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//指定上传文件所存储的路径
	filespath := r.FormValue("filespath")
	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		renderErrors(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	//读取文件字节数
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderErrors(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	//获取上传的文件类型
	detectedFileType := http.DetectContentType(fileBytes)
	filename := filepath.Join(filespath, fileHeader.Filename)
	fmt.Printf("FileType: %s, File: %s\n", detectedFileType, filename)
	//检查目录是否存在，不存在则创建
	if !checkDirIsExist(filespath) {
		if err := os.MkdirAll(filespath, 0755); err != nil {
			renderErrors(w, "UPLOADPATH_NOT_EXISTS", http.StatusInternalServerError)
			return
		}
	}
	//写入文件
	newFile, err := os.Create(filename)
	if err != nil {
		renderErrors(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderErrors(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	//如果文件类型是.zip类型，则解压缩文件到指定目录
	if detectedFileType == "application/zip" {
		//指定解压缩的文件目录
		unzipfilespath := r.FormValue("unzipfilespath")
		if err := unzipfiles(filename, unzipfilespath); err != nil {
			log.Println("unzip file fail,", err)
			w.Write([]byte("fail"))
			return
		}
	}
	w.Write([]byte("SUCCESS"))
}

/*以下为工具方法*/
func renderErrors(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}
func checkDirIsExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//param1:压缩的.zip文件;param2:压缩后所保存的路径
func zipfiles(zipFile, fileDir string) error {
	//创建 zip 包文件
	fw, err := os.Create(zipFile)
	if err != nil {
		log.Fatal()
	}
	defer fw.Close()
	//实例化新的zip.Writer
	zw := zip.NewWriter(fw)
	defer func() {
		//检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	fs, _ := ioutil.ReadDir(fileDir)
	for _, file := range fs {
		fr, err := os.Open(filepath.Join(fileDir, file.Name()))
		if err != nil {
			return err
		}
		fi, err := fr.Stat()
		if err != nil {
			return err
		}
		//写入文件的头信息
		fh, err := zip.FileInfoHeader(fi)
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		//写入文件内容
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return nil
}
func unzipfiles(zipFile, unzipDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := filepath.Join(unzipDir, file.Name)
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}
func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}
func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		panic("start is wrong")
	}
	if end < start || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}
