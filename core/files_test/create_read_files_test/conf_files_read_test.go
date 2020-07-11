package modify_files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

/*
   示例：读取.conf配置文件中的key=value对 add by syf 2020.5.3
   场景：读取jdchain账本配置文件ledger_binding.conf配置文件中的key=val对，定位该参与者的密钥及用户地址信息。
*/
//测试从.conf配置文件中读取key=val对
func TestReadFromConfFile(t *testing.T) {
	//设置读取的文件名称
	filename := "./files/ledger-binding.conf"
	ledgerhash := "j5hAns49vzM6biNXz4AwHh9WjzfhmwuPjZokM5PALFt8mh"
	//读取.conf文件中的k,v对
	conf := readKVConfFiles(filename)
	for k, v := range conf {
		//logs.Info("key=%s ; val=%s",k,v)
		//根据key值获取value值
		if k == strings.Join([]string{"binding", ledgerhash, "parti.address"}, ".") {
			fmt.Println("userAddr:", v)
		} else if k == strings.Join([]string{"binding", ledgerhash, "parti.pk"}, ".") {
			fmt.Println("privkey:", v)
		} else if k == strings.Join([]string{"binding", ledgerhash, "parti.pwd"}, ".") {
			fmt.Println("pwd:", v)
		} else if k == strings.Join([]string{"binding", ledgerhash, "db", "uri"}, ".") {
			fmt.Println("db.uri:", v)
			//截取rocksb的序号
			lastindex := strings.LastIndex(v, "/")
			fmt.Println("lastindex:", lastindex) //51

			db := v[lastindex+1 : len(v)-3]
			fmt.Println("db:", db)

			dbnum := strings.ReplaceAll(db, "rocksdb", "")
			fmt.Println("dbnum:", dbnum)
		}
	}
}

//读取key=value的.conf配置文件
func readKVConfFiles(path string) map[string]string {
	config := make(map[string]string)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}
