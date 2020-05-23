package http_connect

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
   示例：普通http连接Post请求 add by syf 2020.5.19
*/

//测试http连接发送post请求
func TestHttpConnectSendPost(t *testing.T) {

}

func HttpConnectSendPost(url string, param interface{}) ([]byte, error) {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		logs.Error("json格式化数据错误", err)
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonParam))
	if err != nil {
		logs.Error("发送http连接发生错误###", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		logs.Info("dial url:%s fail:%s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logs.Info("http连接请求结果:%s", string(body))
	return body, nil
}
