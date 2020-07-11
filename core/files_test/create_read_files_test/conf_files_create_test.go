package modify_files_test

import (
	"bufio"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

//测试读取ledger.init文件并写入数据到gateway.conf
func TestReadLedgerInitAndWriteGatewayconf(t *testing.T) {
	//获取指定目录下的所有目录(如：network/peer0baas1,peer0baas2)
	filedir := "./files"
	filedirList, err := GetDirList(filedir)
	if err != nil {
		logs.Error("获取目录下的所有目录失败,", err)
		return
	}
	var peer0dir string
	var servername string
	for _, dir := range filedirList {
		//截取最后一个字符串,如:peer0baas1
		dirArr := strings.Split(dir, "/")
		if len(strings.Split(dir, "/")) > 0 {
			servername = dirArr[len(dirArr)-1]
			if strings.Contains(servername, "peer0") {
				peer0dir = dir
				break
			}
		}
	}
	logs.Info("get peer0dir:%s,servername:%s", peer0dir, servername)
	//加载local.conf获取密钥
	localconf := filepath.Join(peer0dir, "local.conf")
	var secretkey []string
	if secretkey = readLocalConfGetSecretKey(localconf); len(secretkey) != 3 {
		logs.Error("read filename:%s get secret key fail...", localconf)
		return
	}
	//加载ledger.init获取参与方个数及host ip
	ledgerinit := filepath.Join(peer0dir, "ledger.init")
	var participant_host []string
	if participant_host = readLedgerInitGetParticipant(ledgerinit); len(participant_host) == 0 {
		logs.Error("read filename:%s get all participant host fail...", localconf)
		return
	}
	//加载gateway.conf,配置参数
	gatewayconf := filepath.Join(peer0dir, "gateway.conf")
	gatewayconf_goal := filepath.Join("/tmp", "gateway.conf")
	if err := modifyGatewayConf(gatewayconf, servername, gatewayconf_goal, secretkey, participant_host); err != nil {
		logs.Error("modify gateway.conf fail,", err)
		return
	}
}

//读取local.conf获取公，私钥,密码
func readLocalConfGetSecretKey(filename string) []string {
	//读取目录下ledger.init配置信息，并将数据写入到gateway.conf
	localconf := filepath.Join(filename)
	cfg_localconf, err := goconfig.LoadConfigFile(localconf)
	if err != nil {
		logs.Error("load local.conf config file fail,err:", err.Error())
		return nil
	}
	if cfg_localconf == nil {
		logs.Error("load lcal.conf.init config is null!")
		return nil
	}
	pubkey, err := cfg_localconf.GetValue("", "local.parti.pubkey")
	if err != nil {
		logs.Error("read filename:%s pubkey fail,", err)
		return nil
	}
	privkey, _ := cfg_localconf.GetValue("", "local.parti.privkey")
	if err != nil {
		logs.Error("read filename:%s privkey fail,", err)
		return nil
	}
	pwd, _ := cfg_localconf.GetValue("", "local.parti.pwd")
	if err != nil {
		logs.Error("read filename:%s pwd fail,", err)
		return nil
	}
	logs.Info("read filename:%s get pubkey:%s,privkey:%s,pwd:%s", pubkey, privkey, pwd)
	if pubkey == "" || privkey == "" || pwd == "" {
		return nil
	}
	var secretkey []string
	secretkey = append(secretkey, pubkey, privkey, pwd)
	return secretkey
}

//读取ledger.init获取所有参与方主机ip
func readLedgerInitGetParticipant(filename string) []string {
	cfg_ledgerinit, err := goconfig.LoadConfigFile(filename)
	if err != nil {
		logs.Error("load ledger.init config file fail,err:", err.Error())
		return nil
	}
	if cfg_ledgerinit == nil {
		logs.Error("load ledger.init config is null!")
		return nil
	}
	ledgerinit_key := cfg_ledgerinit.GetKeyList("DEFAULT")
	ledgerinit_host := []string{}
	for _, v := range ledgerinit_key {
		if strings.Contains(v, ".consensus.host") {
			hostip, err := cfg_ledgerinit.GetValue("", v)
			if err != nil {
				logs.Error("read filename:%s get host ip fail,", err)
				return nil
			}
			ledgerinit_host = append(ledgerinit_host, hostip)
		}
	}
	return ledgerinit_host
}

//配置gateway.conf
func modifyGatewayConf(filename, servername, gatewaydir string, secretkey, hostip []string) error {
	cfg_gateway, err := goconfig.LoadConfigFile(filename)
	if err != nil {
		logs.Error("load gateway.conf fail,err:", err.Error())
		return err
	}
	if cfg_gateway == nil {
		logs.Error("load cfg_gateway is null!")
		return err
	}
	if len(secretkey) == 3 { //设置公私钥及密码
		cfg_gateway.SetValue("", "keys.default.pubkey", secretkey[0])
		cfg_gateway.SetValue("", "keys.default.privkey", secretkey[1])
		cfg_gateway.SetValue("", "keys.default.privkey-password", secretkey[2])
	}
	cfg_gateway.SetValue("", "peer.host", servername)
	cfg_gateway.SetValue("", "peer.size", strconv.Itoa(len(hostip)))
	for k, ip := range hostip {
		fmt.Println("hostip:", hostip)
		cfg_gateway.SetValue("", "peer."+strconv.Itoa(k)+".host", ip)
		cfg_gateway.SetValue("", "peer."+strconv.Itoa(k)+".port", "7080")
		cfg_gateway.SetValue("", "peer."+strconv.Itoa(k)+".secure", "false")
	}
	//注意此处存储gateway.conf文件可选择使用原路径filename，也可使用新路径gatewaydir
	//if err = goconfig.SaveConfigFile(cfg_gateway, filename); err != nil {
	if err = goconfig.SaveConfigFile(cfg_gateway, gatewaydir); err != nil {
		logs.Error("modify gateway.conf config fail,err:", err.Error())
		return err
	}
	logs.Debug("modify gateway.conf success!!!")
	return nil
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

func GetDirList(dirpath string) ([]string, error) {
	var dir_list []string
	dir_err := filepath.Walk(dirpath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				if path != dirpath {
					dir_list = append(dir_list, path)
				}
				return nil
			}

			return nil
		})
	return dir_list, dir_err
}
