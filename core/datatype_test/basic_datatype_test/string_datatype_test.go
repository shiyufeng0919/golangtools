package basic_datatype_test

import (
	"fmt"
	"strings"
	"testing"
)

/*
  示例：string类型 常用操作  add by syf 2020.5.20
*/

//测试去掉某字符串后缀子串的字符串
func TestTrimSuffix(t *testing.T) {
	str := "/usr/local/rocksdb2.db.zip"
	newstr := strings.TrimSuffix(str, ".zip")
	fmt.Printf("deal str:%s delete suffix(.zip) result:%s", str, newstr) //rocksdb2.db
}

//测试某一string是否包含某字符
func TestContainsStr(t *testing.T) {
	str := "/usr/local/peer/rocksdb2.db.zip"
	if strings.Contains(str, "peer") {
		fmt.Printf("str:%s contains peer", str)
		return
	}
	fmt.Printf("str:%s no contains peer", str)
}
