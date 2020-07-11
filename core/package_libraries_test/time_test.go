package package_example_test

import (
	"fmt"
	"testing"
	"time"
)

/*
   示例: golang time包常见用法
*/

//当前时间格式化输出
func TestCurrentTimeFormat(t *testing.T) {
	//注意此处必须为2006-01-02 15:04:05；也可为:2006/01/02 15:04:05
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("currentTime:", currentTime) //2020-06-11 19:27:42
}
