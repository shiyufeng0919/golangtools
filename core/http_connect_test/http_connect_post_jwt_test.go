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
  示例：http连接，携带Token add by syf 2020.5.19
       设置Header  key = Authorization ; Value = Token
*/
//测试http连接 | POST请求 | 携带Token
func TestHttpPostWithToken(t *testing.T) {
	requesturl := "http://127.0.0.1:8080/user"
	token := "eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJqZGNoYWludG9vbHMiLCJzdWIiOiJqZGNoYWludG9vbHMiLCJyb2xlcyI6Im1lbWJlciIsInVzZXJuYW1lIjoiamRjaGFpbnRvb2xzIiwicGFzc3dvcmQiOiJqZGNoYWluIyYjdG9vbHMiLCJpYXQiOjE1ODk4NTU5MjksImV4cCI6MTU4OTg1NzcyOX0.aN4jrDvXrGvlzSrZ82zick4HvaMmR0pCepQRmwADlsI"
	requestparams := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	body, err := HttpPostWithToken(requesturl, requestparams, token)
	if err != nil {
		logs.Error("dial requesturl fail,", err)
		return
	}
	logs.Info("dial requesturl response:", string(body))
}
func HttpPostWithToken(url string, param interface{}, token string) ([]byte, error) {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		logs.Error("json marshal fail,", err)
		return nil, err
	}
	//通过设置tls.Config的InsecureSkipVerify为true，client将不再对服务端的证书进行校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonParam))
	if err != nil {
		logs.Error("http connect fail,", err)
	}
	//设置请求头
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+token)
	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		logs.Info("send request fail,", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logs.Info("dial url:%s response:%s", url, string(body))
	return body, nil
}
