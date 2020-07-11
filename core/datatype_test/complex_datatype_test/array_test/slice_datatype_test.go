package array_test

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestSlice(t *testing.T) {
	arr := []string{}
	arr = append(arr, "A", "B", "C")
	//1,清空数组
	arr = []string{}
	fmt.Println("arr:", arr)
}

func TestAA(t *testing.T) {
	for i := 0; i < 5; i++ {
		logs.Info("I:", i)
		if i != 3 {
			continue
		}
		break
	}
}
