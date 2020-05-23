package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var srv *http.Server

func main() {
	http.HandleFunc("/isReady", isReady)             //准备工作
	http.HandleFunc("/upload", uploadFileHandler)    //上传文件
	http.HandleFunc("/shutdown", shutdownHttpServer) //关闭下载服务
	http.ListenAndServe(":8001", nil)                //设置服务监听端口,下述还有一个文件下载服务端口为8000
}

//准备工作
func isReady(w http.ResponseWriter, r *http.Request) {
	log.Println("isReady...")
	zipDir := "/tmp/zipserver" //打成.zip包所存储的路径及下载服务所指向的目录位置
	//zipFileName:="rocksdb0.zip"
	//zipFileName:="/tmp/peer0/rocks.db/rocksdb0.db"
	filespath := r.FormValue("filespath")
	zipFileName := r.FormValue("dirname") + ".zip"
	log.Printf("zipDir:%s,filespath:%s,zipFileName:%s", zipDir, filespath, zipFileName)
	if !checkDirIsExists(zipDir) {
		if err := os.MkdirAll(zipDir, 0755); err != nil {
			renderError(w, "UPLOADPATH_NOT_EXISTS", http.StatusInternalServerError)
			return
		}
	}
	//压缩rocksdb
	if err := zipfiles(filepath.Join(zipDir, zipFileName), filespath); err != nil {
		log.Printf("zip files error:", err)
		w.Write([]byte("fail"))
	}
	//开启下载服务监听
	if srv == nil || srv.Addr == "" {
		log.Println("ready start download server...")
		go startHttpServer(zipDir)
	}
	log.Println("ready ok...")
	w.Write([]byte("ok"))
}

//启动监听服务，供客户端下载文件
func startHttpServer(monitorDir string) {
	log.Println("startHttpServer...monitorDir:", monitorDir)
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
	log.Println("uploadFileHandler...")
	//指定上传文件所存储的路径
	filespath := r.FormValue("filespath")
	file, fileHeader, err := r.FormFile("uploadFile")
	log.Printf("filespath:%s,filename:%s,size:%s", filespath, fileHeader.Filename, fileHeader.Size)
	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	//读取文件字节数
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("ioutil.ReadAll(file) fail:", err)
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	//获取上传的文件类型
	detectedFileType := http.DetectContentType(fileBytes)
	filename := filepath.Join(filespath, fileHeader.Filename)
	log.Printf("FileType: %s, File: %s\n", detectedFileType, filename)
	//检查目录是否存在，不存在则创建
	if !checkDirIsExists(filespath) {
		if err := os.MkdirAll(filespath, 0755); err != nil {
			log.Printf("mkdir fail,", err)
			renderError(w, "UPLOADPATH_NOT_EXISTS", http.StatusInternalServerError)
			return
		}
	}
	//写入文件
	newFile, err := os.Create(filename)
	if err != nil {
		log.Printf("os.Create(filename) fail:", err)
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		log.Printf("newFile.Write(fileBytes) fail,", err)
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	//解压缩文件
	if detectedFileType == "application/zip" {
		//指定解压缩的文件目录
		unzipfilespath := r.FormValue("unzipfilespath")
		log.Println("unzip files path:", unzipfilespath)
		//param1: 目标文件  param2:解压后的/路径/文件名
		if err := unzipfiles(filename, unzipfilespath); err != nil {
			log.Println("unzip file fail,", err)
			w.Write([]byte("fail"))
			return
		}
	}
	log.Println("uploadFile success...")
	w.Write([]byte("SUCCESS"))
}

/*以下为工具方法*/
func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}
func checkDirIsExists(name string) bool {
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
		log.Printf("zipfiles os.Create(zipFile) fail:", err)
		return err
	}
	defer fw.Close()
	//实例化新的zip.Writer
	zw := zip.NewWriter(fw)
	defer func() {
		//检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Printf("zw.Close() fail:", err)
			return
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
		//file.Name为.zip包中具体文件
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
