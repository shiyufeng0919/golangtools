package transport_files_http

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

/*
  示例：模拟客户端(如浏览器)，调用http连接地址下载文件 add by syf 2020.5.11
  正常浏览器访问:http://127.0.0.1:8080/files/rocksdb.zip即可下载文件
  本示例模拟客户端调用该http地址下载文件

  描述：upload_download_files1_test.go启动到虚拟机node1,并通过vbox将8080端口映射出来
       本机mac请求/upload可上传文件，并可请求http://127.0.0.1:8080/files/rocksdb.zip下载文件
*/

//测试从http服务指定某文件名下载某一文件
func TestDownloadFilesFromHttp(t *testing.T) {
	//指定文件名下载
	requesturl := "http://127.0.0.1:8080/files/rocksdb.zip"
	res, err := http.Get(requesturl)
	if err != nil {
		panic(err)
	}
	//设置文件存储的路径
	savepath := "/tmp/test"
	if !checkDirIsExists(savepath) {
		if err := os.MkdirAll(savepath, 0755); err != nil {
			logs.Error("create dir fail,", err)
			return
		}
	}
	//保存文件
	f, err := os.Create(filepath.Join(savepath, "rocksdb.zip"))
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
}

//测试下载目录中的所有文件
type result struct {
	Args    string            `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

func TestDownloadMultiFilesFromHttp(t *testing.T) {
	//指定文件名下载
	requesturl := "http://127.0.0.1:8080/files"
	resp, err := http.Get(requesturl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
