package struct_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

/*
  示例：struct结构体转Map add by syf 2020.5.10
*/
type userInfo struct {
	Id        int64
	Username  string
	Password  string
	Logintime time.Time
}

//测试struct转map,注意。返回值类型为interface
func TestStructToMap(t *testing.T) {
	userinfo := userInfo{5, "syf", "syf123", time.Now()}
	data := struct2Map(userinfo)
	fmt.Println(data)
}

//struct转map
func struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
