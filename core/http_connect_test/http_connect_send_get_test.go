package http_connect

import (
	"crypto/tls"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
   示例：普通http连接Get请求 add by syf 2020.5.19
*/

//测试http连接发送get请求
func TestHttpConnectSendGet(t *testing.T) {
	requesturl := "http://127.0.0.1:8080/user"
	body, err := HttpConnectSendGet(requesturl)
	if err != nil {
		logs.Error("dial url:%s fail:%s", requesturl, err)
		return
	}
	logs.Info("dial url:%s response:%s", requesturl, string(body))
}

func HttpConnectSendGet(url string) ([]byte, error) {
	//通过设置tls.Config的InsecureSkipVerify为true，client将不再对服务端的证书进行校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//发送http请求
	resp, err := client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logs.Info("dial url:%s response:%s", url, string(body))
	return body, nil
}
