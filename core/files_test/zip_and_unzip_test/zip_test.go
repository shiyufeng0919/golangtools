package zip_and_unzip_test

import (
	"archive/zip"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

/*
  压缩.zip文件 add by syf 2020.5.1
  应用golang原生archive/zip包压缩及解压缩.zip文件
  参考:https://www.jianshu.com/p/aa72b4d50d8d
      https://studygolang.com/articles/7471
*/

//示例：测试压缩目录为.zip文件
func TestZipDir(t *testing.T) {
	//zip名称及存储目录:rocksdb0.zip;zip目标打.zip包的目录
	err := Zip("/tmp/rocksdb0.zip", "/tmp/peer0/rocks.db/rocksdb0.db")
	if err != nil {
		logs.Error("zip file fail,", err)
		return
	}
}

//压缩文件
func Zip(zipFile string, fileDir string) error {
	//创建 zip 包文件
	fw, err := os.Create(zipFile)
	if err != nil {
		log.Fatal()
	}
	defer fw.Close()
	//实例化新的zip.Writer
	zw := zip.NewWriter(fw)
	defer func() {
		//检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	fs, _ := ioutil.ReadDir(fileDir)
	for _, file := range fs {
		fr, err := os.Open(filepath.Join(fileDir, file.Name()))
		if err != nil {
			return err
		}
		fi, err := fr.Stat()
		if err != nil {
			return err
		}
		//写入文件的头信息
		fh, err := zip.FileInfoHeader(fi)
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		//写入文件内容
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return nil
}
