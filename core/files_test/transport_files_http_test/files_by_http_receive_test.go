package transport_files_http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

/*
  示例：模拟服务端接收多文件(接收自files_by_http_upload_test传递的多文件) add by syf 2020.5.8
*/

//测试接收多个文件
func TestReceivesMultiFiles(t *testing.T) {
	http.HandleFunc("/uploadMultiFiles", uploadMultiFilesHandler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}

//上传多文件
func uploadMultiFilesHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//设置内存大小
		r.ParseMultipartForm(32 << 20)
		//获取上传的文件组
		files := r.MultipartForm.File["uploadFile"]
		for i := 0; i < len(files); i++ {
			//打开上传文件
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			//创建上传目录
			os.Mkdir("/tmp/multi", os.ModePerm)
			//创建上传文件
			cur, err := os.Create(filepath.Join("/tmp/multi", files[i].Filename))
			defer cur.Close()
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(cur, file)
			fmt.Println(files[i].Filename) //输出上传的文件名
		}
	})
}
