package basic_datatype

import (
	"fmt"
	"strconv"
	"testing"
)

/*
   示例：float类型常用操作 add by syf 2020.5.20
*/
//测试 string与int间相互转换
func Test_string_int(t *testing.T) {
	//测试： string -> int
	intVal, _ := strconv.Atoi("32000")
	fmt.Println("string convert to int:", intVal)
	//测试: int -> string
	valStr := strconv.Itoa(intVal)
	fmt.Println("int convert string:", valStr)
}

//测试 float与string间相互转换
func Test_Float_String(t *testing.T) {
	//测试： string 转 float64
	float64, _ := strconv.ParseFloat("32000", 64)
	fmt.Println("string conver float64 :", float64)

	//测试: float64 转 string
	float64Str := strconv.FormatFloat(float64, 'f', -1, 64)
	fmt.Println("float64 convert string :", float64Str)
}
