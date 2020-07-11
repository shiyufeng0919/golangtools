package basic_datatype

import (
	"fmt"
	"math"
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
	testfloat64, _ := strconv.ParseFloat("32000", 64)
	fmt.Println("string conver float64 :", testfloat64)

	//测试: float64 转 string
	float64Str := strconv.FormatFloat(testfloat64, 'f', -1, 64)
	fmt.Println("float64 convert string :", float64Str)

	str := "30000"
	float64str, _ := strconv.ParseFloat(str, 64)
	fmt.Println("float64str:", float64str+float64(1000))

}

//测试 float与int间相互转换
func Test_Float_int(t *testing.T) {
	intval := 88
	goalval := float64(intval)
	fmt.Println(goalval)

}

//将float64转成精确的int64
func Wrap(num float64, retain int) int64 {
	return int64(num * math.Pow10(retain))
}

//将int64恢复成正常的float64
func Unwrap(num int64, retain int) float64 {
	return float64(num) / math.Pow10(retain)
}

//精准float64
func WrapToFloat64(num float64, retain int) float64 {
	return num * math.Pow10(retain)
}

//精准int64
func UnwrapToInt64(num int64, retain int) int64 {
	return int64(Unwrap(num, retain))
}
