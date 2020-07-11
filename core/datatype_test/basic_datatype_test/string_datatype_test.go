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

	fmt.Println("=========")
	str2 := ",a,b,c,"
	newstr = strings.Trim(str2, ",")
	fmt.Println("newstr=", newstr)
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

//测试某一字符串头包含某字符
func TestHeadContainsStr(t *testing.T) {
	dir := "/tmp/peer0baasnet1"
	//截取最后一个字符串peer0baasnet1
	dirarr := strings.Split(dir, "/")
	fmt.Println(dirarr)
	if len(dirarr) > 0 {
		goaldir := dirarr[len(dirarr)-1]
		fmt.Println(goaldir)
		if strings.Contains(goaldir, "peer0") {
			fmt.Println("ok")
		}
	}
}

func TestSplitStr(t *testing.T) {
	str := "j5hAns49vzM6biNXz4AwHh9WjzfhmwuPjZokM5PALFt8mh,1"
	newstr := strings.Split(str, ",")
	fmt.Println(newstr[0], "\n", newstr[1])
}

func TestInterfaceToMapArr(t *testing.T) {
	var data interface{}
	var dataArr []map[string]interface{}

	mapdata1 := make(map[string]interface{})
	mapdata1["key1"] = "kaixin"
	mapdata1["val1"] = 100

	mapdata2 := make(map[string]interface{})
	mapdata2["key1"] = "yufeng"
	mapdata2["val1"] = 101

	dataArr = append(dataArr, mapdata1, mapdata2)

	data = dataArr //赋值

	switch data.(type) { //获取data的数据类型，若为[]map[string]interface {}则执行
	case []map[string]interface{}:
		//将interface转成map
		for _, v := range data.([]map[string]interface{}) {
			key1 := v["key1"].(string)
			val1 := v["val1"].(int)
			fmt.Printf("key1:%s,val1:%v \n", key1, val1)
		}
	}
}

func TestWriteFile(t *testing.T) {
	str := ""
	peerIpArr := strings.Split(str, ",")
	fmt.Println(peerIpArr)
}
