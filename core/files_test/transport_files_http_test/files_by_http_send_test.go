package transport_files_http

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

/*
  示例：模拟【客户端】(同postman/浏览器)以http方式上传多文件 add by syf 2020.5.8
  本示例调用files_by_http_download_test路由/uploadMultiFiles以http方式上传多文件
*/

//测试模拟客户端上传文件，等价于postman或浏览器功能。测试发送多个文件
func TestSendMultiFiles(t *testing.T) {
	//请求url,上传文件(此为upload_download_files1_test启动的服务路由)。首先启动upload_download_files1_test测试文件的服务，再调用此方法测试
	requestUrl := "http://127.0.0.1:8081/uploadMultiFiles"
	//requestUrl := "http://127.0.0.1:8080/upload"
	//设置post请求传递的参数示例
	requestParams := map[string]string{
		"networkName": "test",
		"userName":    "syf",
	}
	//设置待发送的所有文件，支持多文件发送
	sendFileNames := []string{}
	//本示例模拟发送单个文件
	filename1 := "/tmp/test/rocksdb.zip"
	filename2 := "/tmp/test/test.txt"
	sendFileNames = append(sendFileNames, filename1, filename2)
	if err := sendMultiFiles(requestParams, sendFileNames, requestUrl); err != nil {
		log.Fatal("send file fail,", err)
		return
	}
}

//测试模拟客户端上传单个文件(向files_upload_download_by_http_test.go中/upload发送请求)
func TestSendSingleFile(t *testing.T) {
	requesturl := "http://127.0.0.1:8001/upload"
	filename := "/tmp/shared/rocksdb0.zip"
	requestparams := map[string]string{
		"filespath":      "/tmp/zippath",   //zip文件存储路径
		"unzipfilespath": "/tmp/unzippath", //zip文件解压缩路径
	}
	if err := sendSingleFile(requestparams, filename, requesturl); err != nil {
		log.Fatal("send file fail,", err)
		return
	}
}

//发送文件(可发送多文件)公共方法
func sendMultiFiles(params map[string]string, filesname []string, url string) error {
	//创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//创建所需上传的文件，一个一个创建
	for _, filename := range filesname {
		//创建第一个需要上传的文件,filepath.Base获取文件的名称。注意此处设置的fieldname值同upload_download_files1_test读取时的值一致
		fileWriter, err := bodyWriter.CreateFormFile("uploadFile", filepath.Base(filename))
		if err != nil {
			log.Fatal("bodyWriter.CreateFormFile err:", err)
			return err
		}
		//打开文件
		fd, err := os.Open(filename)
		if err != nil {
			log.Fatal("os.Open filename err:", err)
			return err
		}
		defer fd.Close()
		//把第一个文件流写入到缓冲区里去
		_, err = io.Copy(fileWriter, fd)
		if err != nil {
			log.Fatal("io.Copy err:", err)
			return err
		}
	}
	//这一句写入附加字段必须在_,_=io.Copy(fileWriter,fd)后面
	//params := map[string]string{
	//	"networkName": networkName,
	//	"ledgerName":  ledgerName, //初始设置为系统账本
	//}
	if len(params) != 0 {
		//param是一个一维的map结构
		for k, v := range params {
			bodyWriter.WriteField(k, v)
		}
	}
	//获取请求Content-Type类型,后面有用
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	//创建一个http客户端请求对象
	client := &http.Client{}
	log.Println("send files to url:", url)
	//创建一个post请求
	req, _ := http.NewRequest("POST", url, nil)
	//这里的Content-Type值就是上面contentType的值
	req.Header.Set("Content-Type", contentType)
	//转换类型
	req.Body = ioutil.NopCloser(bodyBuf)
	//发送数据
	data, _ := client.Do(req)
	//读取请求返回的数据
	bytes, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal("send jdchaintools request fail,err:", err.Error())
		return err
	}
	defer data.Body.Close()
	log.Println("jdchaintools response:", string(bytes))
	//解析响应结果
	return nil
}

//发送文件(单个文件)公共方法
func sendSingleFile(params map[string]string, filesname string, url string) error {
	//创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//创建第一个需要上传的文件,filepath.Base获取文件的名称。注意此处设置的fieldname值同upload_download_files1_test读取时的值一致
	fileWriter, err := bodyWriter.CreateFormFile("uploadFile", filepath.Base(filesname))
	if err != nil {
		log.Fatal("bodyWriter.CreateFormFile err:", err)
		return err
	}
	//打开文件
	fd, err := os.Open(filesname)
	if err != nil {
		log.Fatal("os.Open filename err:", err)
		return err
	}
	defer fd.Close()
	//把第一个文件流写入到缓冲区里去
	_, err = io.Copy(fileWriter, fd)
	if err != nil {
		log.Fatal("io.Copy err:", err)
		return err
	}
	if len(params) != 0 {
		//param是一个一维的map结构
		for k, v := range params {
			bodyWriter.WriteField(k, v)
		}
	}
	//获取请求Content-Type类型,后面有用
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	//创建一个http客户端请求对象
	client := &http.Client{}
	log.Println("send files to url:", url)
	//创建一个post请求
	req, _ := http.NewRequest("POST", url, nil)
	//这里的Content-Type值就是上面contentType的值
	req.Header.Set("Content-Type", contentType)
	//转换类型
	req.Body = ioutil.NopCloser(bodyBuf)
	//发送数据
	data, _ := client.Do(req)
	//读取请求返回的数据
	bytes, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal("send jdchaintools request fail,err:", err.Error())
		return err
	}
	defer data.Body.Close()
	log.Println("jdchaintools response:", string(bytes))
	return nil
}
