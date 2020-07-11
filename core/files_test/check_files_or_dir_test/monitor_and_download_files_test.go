package check_files_or_dir

import (
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

/*
   示例：监控远端服务某文件是否存在，若存在则下载，否则一直监测 add by syf 2020.6.10
   项目需求：jdchaintools定时监测peer节点下启动的(main.go)服务，上传local.conf及ledger_binding.conf配置文件到jdchaintools指定目录下。
*/

//监控远端是否存在某文件，存在，则下载，否则一直监测. 对于local.conf文件的监测(peer启动生成local.conf)及下载
func TestMonitorAndDownloadFiles(t *testing.T) {
	//指定服务下载路径
	downloadPath := `http://127.0.0.1:8001/files`
	for i := 0; i < 120; i++ { //尝试10min
		time.Sleep(time.Second * 5)
		//检测http下载服务是否已启动
		if checkHttpServerIsAvailable(downloadPath) {
			break
		}
	}
	//指定文件名下载,监测文件是否已存在，若不存在，则始终监测
	filename := `local.conf`
	//注意：此处不能用filepath.join,否则http://后边的/有一个会被当作转义字符
	requesturl := downloadPath + "/" + filename
	logs.Info("ready dial requesturl=", requesturl)
	var res *http.Response
	var err error
	for i := 0; i < 120; i++ { //重试10min
		time.Sleep(time.Second * 5) //每5s重试一次
		res, err = http.Get(requesturl)
		if err != nil {
			logs.Error("download files from url:%s fail:%s")
			continue
		}
		logs.Info("download files from url:%s response statusCode:%s,status:%s,body:%s", requesturl, res.StatusCode, res.Status, res.Body)
		if res.StatusCode == http.StatusOK && res.Body != nil {
			logs.Info("sucess download local.conf...")
			break
		} else {
			logs.Info("no generate local.conf...")
			continue
		}
	}
	//设置文件存储的路径
	savepath := "/tmp/test"
	if !checkDirIsExists(savepath) {
		if err := os.MkdirAll(savepath, 0755); err != nil {
			logs.Error("create dir fail,", err)
			return
		}
	}
	//保存文件
	f, err := os.Create(filepath.Join(savepath, "local.conf"))
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
}

//监控远端是否存在某文件，存在，则下载，否则一直监测. 同上逻辑。但获得到ledger-binding.conf文件后需检查是否与本地文件一致，若不一致，则保存。否则继续监测。
//ledger-binding.conf为创建账本所生成，每生成一个账本，此配置文件会追加一个账本
func TestMonitorAndDownloadLedgerbinding(t *testing.T) {
	//指定服务下载路径
	downloadPath := `http://127.0.0.1:8001/files`
	for i := 0; i < 120; i++ { //尝试10min
		time.Sleep(time.Second * 5)
		//检测http下载服务是否已启动
		if checkHttpServerIsAvailable(downloadPath) {
			break
		}
	}
	//指定文件名下载,监测文件是否已存在，若不存在，则始终监测
	filename := `ledger-binding.conf`
	//注意：此处不能用filepath.join,否则http://后边的/有一个会被当作转义字符
	requesturl := downloadPath + "/" + filename
	logs.Info("ready dial requesturl=", requesturl)
	var res *http.Response
	var err error
	for i := 0; i < 180; i++ { //重试15min
		time.Sleep(time.Second * 5) //每10s重试一次
		res, err = http.Get(requesturl)
		if err != nil {
			logs.Error("download files from url:%s fail:%s")
			continue
		}
		logs.Info("download files from url:%s response statusCode:%s,status:%s,body:%s", requesturl, res.StatusCode, res.Status, res.Body)
		if res.StatusCode == http.StatusOK && res.Body != nil {
			logs.Info("sucess download ledger-binding.conf...")
			//设置文件存储的路径
			savepath := "/tmp/test"
			if !checkDirIsExists(savepath) {
				if err := os.MkdirAll(savepath, 0755); err != nil {
					logs.Error("create dir fail,", err)
					return
				}
			}
			//检查文件是否存在
			goalfile := filepath.Join(savepath, filename)
			//是否保存文件标识
			saveFiles := false
			//若文件不存在，则直接写入到指定目录即可
			if !FileIsExists(goalfile) {
				saveFiles = true
			} else { //检测文件是否一致，若文件不一致，则保存文件，否则再次请求下载ledger-binding.conf
				srcFilesSize := getFileSize(goalfile)
				logs.Info("old files size:", srcFilesSize)
				newFilesSize := res.ContentLength
				logs.Info("download files size:", newFilesSize)
				if srcFilesSize < newFilesSize {
					saveFiles = true
				}
			}
			//保存文件
			if saveFiles {
				f, err := os.Create(goalfile)
				if err != nil {
					panic(err)
				}
				io.Copy(f, res.Body)
				logs.Info("save ledger_binding.conf success...")
				break
			}
		} else {
			logs.Info("no generate ledger-binding.conf...")
			continue
		}
	}
}

//检查http server是否可用 add by syf 2020.6.10
func checkHttpServerIsAvailable(goalurl string) bool {
	logs.Info("ready check goalurl:%s is available...", goalurl)
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(goalurl)
	if err == nil && resp.StatusCode == http.StatusOK {
		logs.Info("goalurl:%s is available...", goalurl)
		return true
	}
	return false
}

//检查目录是否存在
func checkDirIsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//检查文件是否存在
func FileIsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//获取文件大小
func getFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
