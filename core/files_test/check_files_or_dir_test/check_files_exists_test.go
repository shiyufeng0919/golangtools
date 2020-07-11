package check_files_or_dir_test

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

/*读写文件方式：1,ioutil 2,os 3,bufio 其中ioutil效率更高些 add by syf 2020.5.31*/

//示例1：检查某目录下的文件是否存在 add by syf 2020.5.22
func TestCheckFilsIsExsits(t *testing.T) {
	filenames := "/tmp/shared/ledger-binding.conf"
	if CheckFilesIsExists(filenames) {
		fmt.Println("文件存在...")
		return
	}
	fmt.Println("文件不存在...")
}

//示例2: 向文件写入数据，文件所在路径及文件没有则创建。 add by syf 2020.5.31
const (
	NewDirMode  = 0755
	NewFileMode = 0755
)

func TestIoUtilSaveFiles(t *testing.T) {
	filename := "/tmp/haha/syf.txt"
	writedata := "越努力，越幸运! fighting..."
	if err := SaveFile(filename, writedata); err != nil {
		logs.Error("save file fail,", err)
		return
	}
}

//示例3:读取文件
func TestIoutilReadFiles(t *testing.T) {
	filename := "/tmp/syf.txt"
	readcontent, _ := ioutil.ReadFile(filename)
	fmt.Println("读取内容：", string(readcontent))
}

//示例4: 读取文件
func TestOsReadFiles(t *testing.T) {
	filename := "/tmp/syf.txt"
	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	contentByte, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	fmt.Println("读取内容：", string(contentByte))
}

func SaveFile(fpath string, data interface{}) error {
	logs.Debug("save to file: %s", fpath)
	var valueBytes []byte
	var err error
	switch data.(type) {
	case string:
		valueBytes = []byte((data).(string))
	case []byte:
		valueBytes = (data).([]byte)
	default:
		if valueBytes, err = json.Marshal(data); err != nil {
			return err
		}
	}
	dir := filepath.Dir(fpath)
	if !FileExists(dir) {
		if err = os.MkdirAll(dir, NewDirMode); err != nil {
			logs.Error("fail to make directory: %s, with error: %v", dir, err)
			return err
		}
	}
	if err = ioutil.WriteFile(fpath, valueBytes, NewFileMode); err != nil {
		logs.Error("fail to save file: %s, with error: %v", fpath, err)
		return err
	}
	return nil
}
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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
