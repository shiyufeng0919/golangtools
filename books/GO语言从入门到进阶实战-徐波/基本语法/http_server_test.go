package 基本语法

import (
	"net/http"
	"testing"
)

//模拟http文件服务器,浏览器请求:http://127.0.0.1:8080/ 会显示当前目录"基础语法"下所有文件列表
//知识点：net/http包，这个包的作用是HTTP的基础封装和访问
func TestHttpServer(t *testing.T) {
	//将当前目录在HTTP服务器上映射
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)
}
