package modify_files

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	NewDirMode  = 0755
	NewFileMode = 0755
)

var (
	rootDir  = "/tmp"
	network  = "syfnet"
	orgname  = "syf"
	channel  = "syfchannel"
	peername = []string{"baasnet1peer0", "baasnet1peer1", "baasnet1peer2"}
)

//示例：覆盖写:测试创建.txt文件并向里写入内容方式一
func TestTxtFilesCreate1(t *testing.T) {
	filePath := filepath.Join(rootDir, network, orgname, channel)
	for _, peer := range peername {
		filename := filepath.Join(filePath, peer+".txt")
		writeData1 := "31000"
		writeData2 := "32000"
		writeData := writeData1 + "\n" + writeData2
		SaveFile(filename, writeData)
	}
}

//示例：测试创建.txt文件并向里写入内容方式二
func TestTxtFilesCreate2(t *testing.T) {
	filePath := filepath.Join(rootDir, network, orgname, channel)
	for _, peer := range peername {
		filename := filepath.Join(filePath, peer+".txt")
		writeData1 := "32000"
		writeData2 := "33000"
		writeData := writeData1 + "\n" + writeData2
		if !CheckDirIsExists(filePath) {
			if err := os.MkdirAll(filePath, 0755); err != nil {
				logs.Error("create fileDir:%s fail,%s", filePath, err)
				return
			}
		}
		if err := WriteToFile(filename, writeData); err != nil {
			logs.Error("write data to txt fail,", err)
			return
		}
	}
}

//测试写入数组数据到文件，并将数组元素以逗号分隔写入文件
func TestTxtFilesAppendWrite(t *testing.T) {
	filename := "/tmp/test.txt"
	floatdata := []float64{36004, 36005}
	stringdata := []string{}
	for _, v := range floatdata {
		stringdata = append(stringdata, strconv.FormatFloat(v, 'f', -1, 64))
	}
	if err := WriteArrayToFiles(stringdata, filename); err != nil {
		logs.Error("write data fail,", err)
		return
	}
	fmt.Println("输出最大端口:", floatdata[1], ",端口+1:", floatdata[1]+1)
}

//读取上述追加写入文件的内容
func TestTxtFilesAppendRead(t *testing.T) {
	filename := "/tmp/test.txt"
	//文件存在，再读取
	if checkFileIsExist(filename) {
		val, err := ioutil.ReadFile(filename)
		if err != nil {
			logs.Error("read file fail,", err)
			return
		}
		fmt.Println("读取的值:", string(val))
		//以逗号拆分数据
		valArr := strings.Split(strings.Trim(string(val), ","), ",")
		fmt.Println("最大端口值:", valArr[len(valArr)-1])
		valFloat64, _ := strconv.ParseFloat(valArr[len(valArr)-1], 64)
		fmt.Println("最大端口值+1:", valFloat64+1)
		for _, v := range valArr {
			fmt.Println("端口值:", v)
		}
	}
	logs.Info("file is not exists...")

}

//示例：测试读取.txt文件方式一
func TestTxtFilesRead1(t *testing.T) {
	filePath := filepath.Join(rootDir, network, orgname, channel)
	for _, peer := range peername {
		filename := filepath.Join(filePath, peer+".txt")
		val, _ := ioutil.ReadFile(filename)
		//str:=strings.Split(val,",")
		//fmt.Printf("port1:%s,port2:%s \n",str[0],str[1])
		fmt.Printf("peer:%s,val:%s \n", peer, val)
	}
}

//示例：测试读取.txt文件方式二
func TestTxtFilesRead2(t *testing.T) {
	filePath := filepath.Join(rootDir, network, orgname, channel)
	for _, peer := range peername {
		filename := filepath.Join(filePath, peer+".txt")
		val := ReadFromFile(filename)
		if val == "" {
			break
		}
		//str:=strings.Split(val,",")
		//fmt.Printf("port1:%s,port2:%s \n",str[0],str[1])
		fmt.Printf("peer:%s,val:%s \n", peer, val)
	}
}

func ReadFromFile(filename string) string {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	contentByte, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	//fmt.Println(string(contentByte))
	return string(contentByte)
}

func WriteToFile(fileName, data string) error {
	var f *os.File
	var err error

	if !checkFileIsExist(fileName) { //文件不存在则创建
		logs.Info("file not exsits...")
		f, err = os.Create(fileName) //创建文件
		if err != nil {
			logs.Error("create file fail:%s", err)
			return err
		}
	} else { //文件存在，则直接写入
		logs.Info("file exists....")
		f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			logs.Error("os.open file fail,err:", err)
			return err
		}
	}
	_, err = io.WriteString(f, data) //写入文件(字符串)
	if err != nil {
		logs.Error("data write to file fail,filename:%s,err:%s", fileName, err)
		return err
	}
	return nil
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	logs.Info("checkFileIsExist:", checkFileIsExist)
	return exist
}

//检查目录是否存在
func CheckDirIsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//保存数据到文件
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

//写入数组到文件,追加写,并以逗号分割
func WriteArrayToFiles(m []string, fileName string) error {
	//O_CREATE:创建文件,如果文件不存在 ; O_WRONLY:只写模式 ; O_APPEND:追加内容(若需要追加写，只需添加|os.O_APPEND即可)
	//os.O_TRUNC 覆盖写入，不加则追加写入
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range m {
		lineStr := fmt.Sprintf("%v", v)
		//fmt.Fprintln(w, lineStr)
		fmt.Fprint(w, lineStr, ",")
	}
	return w.Flush()
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
