package check_files_or_dir_test

import (
	"fmt"
	"os"
	"testing"
)

//示例1：检查某目录下的文件是否存在 add by syf 2020.5.22
func TestCheckFilsIsExsits(t *testing.T) {
	filenames := "/tmp/shared/ledger-binding.conf"
	if CheckFilesIsExists(filenames) {
		fmt.Println("文件存在...")
		return
	}
	fmt.Println("文件不存在...")
}

//检查文件是否存在
func CheckFilesIsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
