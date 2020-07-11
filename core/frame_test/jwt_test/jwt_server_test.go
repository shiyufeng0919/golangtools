package jwt_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jwt "github.com/shiyufeng0919/golangtools/core/frame_test/jwt_test"
	"testing"
)

/*
  示例： jwt(json web token)鉴权 add by syf 2020.5.23

  JWT认证：用户注册后, 服务器生成一个 JWT token返回给浏览器, 浏览器向服务器请求数据时将 JWT token 发给服务器, 服务器用 signature 中定义的方式解码
  一个 JWT token包含3部分：
  1，header: 告诉我们使用的算法和token类型 ;
  2,Payload: 必须使用 sub key 来指定用户 ID, 还可以包括其他信息比如 email, username 等.
  3,Signature: 用来保证 JWT 的真实性. 可以使用不同算法

  另：可了解cookie & session & token区别
*/
func TestJwtAuthorize(t *testing.T) {
	//router := JwtRouter() //方式一
	router := JwtRouter2() //方式二
	router.Run(":" + "8000")
}

//初始化jwt路由
func JwtRouter() *gin.Engine {
	router := gin.Default()
	jwtRouterGroup := router.Group("/jwt")
	//初始化jwt认证
	jwt, err := jwt.InitJwt()
	if err != nil {
		panic(fmt.Sprintf("fail to initialize jwt: %#v", err))
	}
	//登录认证，获取token。POST请求http://127.0.0.1:8000/jwt/login,设置Authorization->Type=Basic Auth->输入username和password
	jwtRouterGroup.POST("/login", jwt.LoginHandler)
	//示例：所有以/network组的路由均需jwt鉴权。
	//POST请求http://127.0.0.1:8000/jwt/network/add需设置Authorization->Bearer Token值(上述token)
	networkRouter := jwtRouterGroup.Group("/network").Use(jwt.MiddlewareFunc())
	networkRouter.POST("/add", AddNetwork)
	return router
}

func JwtRouter2() *gin.Engine {
	router := gin.Default()
	jwtRouterGroup := router.Group("/jwt")
	//初始化jwt认证
	jwt, err := jwt.InitJwt()
	if err != nil {
		panic(fmt.Sprintf("fail to initialize jwt: %#v", err))
	}
	//登录认证，获取token。POST请求http://127.0.0.1:8000/jwt/login,设置Authorization->Type=Basic Auth->输入username和password
	jwtRouterGroup.POST("/login", jwt.LoginHandler)
	jwtRouterGroup.Use(jwt.MiddlewareFunc()) //jwtRouterGroup下面所有的路由均生效
	//示例：所有以/network组的路由均需jwt鉴权。
	//POST请求http://127.0.0.1:8000/jwt/network/add需设置Authorization->Bearer Token值(上述token)
	networkRouter := jwtRouterGroup.Group("/network")
	networkRouter.POST("/add", AddNetwork)
	return router
}

//添加网络
func AddNetwork(ctx *gin.Context) {
	fmt.Println("add network...")
}
