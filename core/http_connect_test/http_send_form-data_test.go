package http_connect_test

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
)

/*
	示例: golang发送form-data格式请求参数  add by syf 2020.5.12

*/
//golang发送form-data格式请求,参见示例：https://www.jianshu.com/p/51b0a14429d0
func TestHttpSendFormData(t *testing.T) {
	postData := make(map[string]string)
	postData["networkname"] = "net1"
	postData["filepath"] = "/tmp"
	url := "http://127.0.0.1:8000/upload"
	PostWithFormData("POST", url, &postData)
}

//设置form-data请求参数，解析form-data请求参数见files_test/files_upload_donwload_by_http_test.go
func PostWithFormData(method, url string, postData *map[string]string) ([]byte, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, nil
}
