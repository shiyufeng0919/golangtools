package jwt

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"git.jd.com/baas/entnet/jd-baas-rtmc/pkg"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
	"time"
)

/*
  模拟jwt客户端登录jwt服务/login，获取token add by syf 2020.6.18
*/
type HttpRequestClient struct {
	Scheme       string
	Method       string
	ReqPath      string
	WithLogin    bool
	LoginUser    string
	LoginPwd     string
	ReqHeader    map[string]string
	Timeout      int
	UrlParam     map[string]string
	JsonParam    interface{}
	ReqBody      []byte
	FilePost     map[string]string
	ResponseData interface{}
}
type AgentResponse struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

func TestJwtClient(t *testing.T) {
	router := jwtClientRouter() //方式二
	router.Run(":" + "8001")
}

func jwtClientRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/jwtClient", getToken)
	return router
}

//调用jwt_server_test 获取token
func getToken(ctx *gin.Context) {
	requesturl := "http://127.0.0.1:8000/jwt/login"
	username := "syf"
	userpwd := "kaixin"
	agentResp := &AgentResponse{}
	if err := loginAgent(requesturl, username, userpwd, agentResp); err != nil {
		logs.Error("login agent fail...")
	}
	logs.Info("agent response:", agentResp)
}

//调用jdchainAgent发送/login请求
func loginAgent(inReqAllPath, inUserName, inPwd string, outRespData interface{}) error {
	client := &pkg.HttpRequestClient{}
	client.Method = pkg.ReqMethodPost
	client.ReqPath = inReqAllPath
	client.ReqHeader = make(map[string]string)
	//client.ReqHeader[pkg.ReqHeaderAuthorization] = pkg.JwtAuthPrefix
	client.ReqHeader[pkg.ReqHeaderConnectType] = pkg.ReqDataJson
	client.WithLogin = true
	client.LoginUser = inUserName
	client.LoginPwd = inPwd
	client.ResponseData = outRespData
	if err := client.SendRequest(); err != nil {
		return err
	}
	return nil
}

//发送请求
func (client *HttpRequestClient) sendRequest() error {
	// 如果未设置读写及连接超时时间,则默认600s,即10分钟
	if client.Timeout == 0 {
		client.Timeout = 600
	}
	ConnTimeout := time.Duration(client.Timeout) * time.Second
	RWTimeout := time.Duration(client.Timeout) * time.Second
	var req *httplib.BeegoHTTPRequest
	switch client.Method {
	case "POST":
		req = httplib.Post(client.ReqPath)
	case "GET":
		req = httplib.Get(client.ReqPath)
	case "PUT":
		req = httplib.Put(client.ReqPath)
	case "DELETE":
		req = httplib.Delete(client.ReqPath)
	default:
		logs.Info("Unexpected request method")
		return errors.New("Unexpected request method")
	}
	// 访问对象是否需要登录
	if client.WithLogin {
		req.SetBasicAuth(client.LoginUser, client.LoginPwd)
	}
	// 设置连接和读写超时时间
	req = req.SetTimeout(ConnTimeout, RWTimeout)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // 设置tls为跳过
	for k, v := range client.ReqHeader {
		req.Header(k, v)
	}
	for k, v := range client.UrlParam {
		req.Param(k, v)
	}
	for k, v := range client.FilePost {
		req.PostFile(k, v)
	}
	if string(client.ReqBody) != "[]" && client.ReqBody != nil {
		req.Body(client.ReqBody)
	}
	req.JSONBody(client.JsonParam) // 设置请求体参数
	// 发送请求,获取响应内容
	response, err := req.Response()
	if err != nil {
		logs.Error("request send to remote server is fail, err: ", err)
		return err
	}
	// 如果响应的内容不为200,则认为请求出现错误
	if response.StatusCode != http.StatusOK {
		logs.Error("get response code from remote server is err, is not 200: %d, %s", response.StatusCode, response.Status)
		return errors.New(fmt.Sprintf("request remote server is err, status: %s", response.Status))
	}
	// 以字符串的形式获取响应体
	respBodyByte, err := req.Bytes()
	if err != nil {
		logs.Error("get response body from remote server is fail, err: ", err)
		return err
	}
	//logs.Debug("response body from remote server content is: ", string(respBodyByte))
	if string(respBodyByte) == "[]" {
		//logs.Error("response body is null: ", string(respBodyByte), len(respBodyByte), respBodyByte) [common.go:131] response body is null:  [] 2 [91 93]
		logs.Warn("response body is null")
		return nil
	}
	if client.ResponseData == nil { // 如果不需要接收返回值的内容,则直接将client的返回值置空即可
		return nil
	}
	err = json.Unmarshal(respBodyByte, client.ResponseData)
	if err != nil {
		logs.Error("response body to unmarshal err: ", err)
		return err
	}
	return nil
}
