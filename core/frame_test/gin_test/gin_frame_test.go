package gin_test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

/*
   示例: Gin框架应用 add by syf 2020.6.11
   参考Go语言框架gin中文文档:https://github.com/skyhee/gin-doc-cn
*/

//Gin框架
func TestGinFrame(t *testing.T) {
	router := gin.Default()
	http.ListenAndServe(":8001", router)
}
