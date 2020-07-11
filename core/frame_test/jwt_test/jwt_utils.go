package jwt

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
  示例：jwt工具类 add by syf 2020.5.23
*/
const (
	identityKey = "id"
	ROLE_USER   = 1
)

var apiuser LoginUser

func init() {
	apiuser.UserName = "syf"
	apiuser.Password = "kaixin"
}

type LoginUser struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserIdentity struct {
	UserName string
	Role     int
}

type Return struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"err_msg"`
}

//初始化jwt
func InitJwt() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 30,
		MaxRefresh:  time.Minute * 30,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserIdentity); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		//身份
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &UserIdentity{UserName: claims["id"].(string)}
		},
		//登录认证
		Authenticator: func(c *gin.Context) (interface{}, error) {
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				logs.Error("login request less username and password!")
				return "", jwt.ErrMissingLoginValues
			}
			//验证环境变量apiuser&apipwd是否与请求一致
			fmt.Printf("username:%s,pwd:%s", username, password)
			fmt.Printf("apiuser:%s,apipwd:%s", apiuser.UserName, apiuser.Password)
			if username != apiuser.UserName || password != apiuser.Password {
				logs.Error("Authenticator fail!")
				return "", jwt.ErrInvalidAuthHeader
			}
			loginVals := &LoginUser{username, password}
			return &UserIdentity{UserName: loginVals.UserName, Role: ROLE_USER}, nil
		},
		//授权
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*UserIdentity); ok {
				return true
			} else {
				logs.Error("Authorizator fail！")
				return false
			}
		},
		//未授权
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, ReturnUnauthorized(fmt.Sprintf("授权失败!")))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer", //设置token头，可根据实际需求修改此值
		TimeFunc:      time.Now,
	})
}

//返回未授权响应结果
func ReturnUnauthorized(msg string) *Return {
	return &Return{
		Code:   http.StatusUnauthorized,
		Status: "error",
		ErrMsg: msg,
	}
}
