package modify_files_test

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"strconv"
	"testing"
)

/*
  目标：生成.conf格式文件 add by syf 2020.5.3
  项目需求：jdchain账本初始化配置文件ledger_binding.conf,动态入网及动态加入账本需生成该配置文件。
  [文件ledger_binding.conf内容示例：]
  	ledger.bindings=xxx, \ yyy
    binding.xxx.name=first #注意，此中间值为ledger.bindings生成的值
    ...
    binding.yyy.name=second

  本示例实现：1。将数组写入.conf文件(有序) 2。将map写入.conf文件(无序)
*/

//【方法一】：写入数组数据到.conf格式配置文件 (有序)
func TestWriteToConfFile(t *testing.T) {
	filePath := "ledger-binding.conf"
	writeData := []string{}
	ledger := []string{}
	//账本Hash以数组形式模拟
	ledgerHash := []string{"j5wtYVvnzoeRSPQDrMvL4nWfobe3yGQT7tBnQzuLkhS6CV", "j5hAns49vzM6biNXz4AwHh9WjzfhmwuPjZokM5PALFt8mh"}
	allLedger := "" //所有账本，以,\分隔
	for k, hash := range ledgerHash {
		if k == 0 { //第一个
			allLedger = "ledger.bindings=" + hash //第一个账本
		} else { //不是第一个
			allLedger = allLedger + ", \\" + hash //其他账本，中间.\分隔
		}
		ledger = append(ledger, "binding."+hash+".name="+"test"+strconv.Itoa(k))                                                                //账本名称
		ledger = append(ledger, "binding."+hash+".parti.address="+"LdeNuC9Hz75qvXAawDXxYem7ybYW2cgpFb6rm")                                      //参与者用户地址
		ledger = append(ledger, "binding."+hash+".parti.name="+"baidu.com")                                                                     //参与者用户名称
		ledger = append(ledger, "binding."+hash+".parti.pk="+"177gk1KSG3mvYtvxNMZfZmFqB1dJPfQUUArrkHtLm8swfKdeNjVt2mjp6ntcwgWUuSp8m9B")         //参与者私钥
		ledger = append(ledger, "binding."+hash+".parti.pwd="+"8eaDBRNtFyyfoAXLjiR2TZsfMfkiKurTqU95GDm9FL2H")                                   //参与者解密私钥密码
		ledger = append(ledger, "binding."+hash+".db.uri="+"rocksdb:///home/jdchain1.1.4/jdchain/peer0/rocks.db/rocksdb"+strconv.Itoa(k)+".db") //参与者rocksdb目录
	}
	//重新构建数组，设置第0个元素值
	writeData = append(append(writeData, allLedger), ledger...)
	//写入数组数据到文件
	if err := WriteArrayToFiles(writeData, filePath); err != nil {
		logs.Error("arrays write file fail,", err)
		return
	}
	logs.Info("success...")
}

//写入数组数据到文件
func WriteArrayToFiles(m []string, fileName string) error {
	//O_CREATE:创建文件,如果文件不存在 ; O_WRONLY:只写模式 ; O_APPEND:追加内容(若需要追加写，只需添加|os.O_APPEND即可)
	//os.O_TRUNC 覆盖写入，不加则追加写入
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range m {
		lineStr := fmt.Sprintf("%s", v)
		fmt.Fprintln(w, lineStr)
	}
	return w.Flush()
}

//【方法二】：将map[string]string写入文件，问题点：无序
func WriteMaptoFile(m map[string]string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("create map file error: %v\n", err)
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, v := range m {
		lineStr := fmt.Sprintf("%s^%s", k, v)
		fmt.Fprintln(w, lineStr)
	}
	return w.Flush()
}
