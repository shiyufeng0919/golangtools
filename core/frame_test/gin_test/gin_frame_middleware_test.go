package gin

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

/*
   示例: Gin框架中间件项目应用，拦截所有路由，并将请求业务操作记录日志 add by syf 2020.6.11
*/
var routerPaths map[string]string = map[string]string{
	"/test/add":    "添加参数",
	"/test/del":    "删除参数",
	"/test/update": "修改参数",
}

//此处结构体，有可能前端传入参数key值不一致...
type proxyParams struct {
	NetworkName  string `json:"networkName"`
	Network_Name string `json:"network_name"`
	UserName     string `json:"userName"`
}

//中间件
func TestGinFrameMiddleware(t *testing.T) {
	router := gin.Default()
	//router.Use(myMiddleware()) //全局中间件
	routerGroup := router.Group("/test")
	routerGroup.Use(myMiddleware()) //路由组中间件
	{
		routerGroup.POST("/add", addData) //请求路由:http://localhost:8001/test/add 发现问题:127.0.0.1报404
	}
	router.Run(":8001")
}

//自定义中间件
func myMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求路由
		requestRouter := c.Request.URL.Path
		//当前时间
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		logs.Info("request myMiddleware currentTime:", currentTime)
		//路由在维护的map中
		if _, ok := routerPaths[requestRouter]; ok {
			logs.Info("requestRouter %s match,op:%s", requestRouter, routerPaths[requestRouter])
			body, err := c.GetRawData()
			if err != nil {
				logs.Error("request params fail,", err)
				return
			}
			//注意此行必须填加，否则c.Next调用/add路由时，会得不到参数(即此处逻辑是把参数再放回到请求体中...)
			//此处获取完数据后立即写回去。另：存在问题，上传文件会有问题..
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			params := new(proxyParams)
			if err = json.Unmarshal(body, &params); err != nil {
				logs.Error("unmarshal body to params fail,", err)
				return
			}
			networkName := strings.Join([]string{params.NetworkName, params.Network_Name}, "")
			if networkName != "" {
				//从cookie中获取loginName
				loginname, _ := c.Cookie("loginName")
				logs.Info("do save network=%s , username=%s data to log...", networkName, loginname)
			}
			c.Next()
		}
	}
}

//添加数据
func addData(ctx *gin.Context) {
	params := proxyParams{}
	if err := ctx.BindJSON(&params); err != nil {
		logs.Error("addData params is error:", err)
		return
	}
	logs.Info("addData request params network:%s,username:%s", params.NetworkName, params.UserName)
	return
}
