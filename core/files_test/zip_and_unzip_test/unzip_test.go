package zip_and_unzip

import (
	"archive/zip"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
  解压缩.zip文件到指定目录 add by syf 2020.5.1
  应用golang原生archive/zip包压缩及解压缩.zip文件
  参考:https://www.jianshu.com/p/aa72b4d50d8d
      https://studygolang.com/articles/7471
*/

var zipFile = "/tmp/rocksdb0.zip"     //目标：待解压的.zip文件
var unzipDir = "/tmp/shared/rocksdb0" //解压缩后文件存储目录
//示例：测试解压缩目录
func TestUnzipDir(t *testing.T) {
	//参数:目标zip文件；指定解压缩目录
	err := Unzip(zipFile, unzipDir)
	if err != nil {
		logs.Error("unzip dir fail,", err)
		return
	}
}

//解压缩，参数:目标.zip文件；解压缩存储目录
func Unzip(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := filepath.Join(dest, strings.TrimSuffix(file.Name, ".zip"))
		fmt.Println("unzip filename:", filename)
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}
func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}
func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		panic("start is wrong")
	}
	if end < start || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}
