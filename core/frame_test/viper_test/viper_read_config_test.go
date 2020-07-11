package viper_test

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"testing"
)

/*
	示例: viper读取配置文件 add by syf 2020.6.11
	github.com/spf13/viper go viper包介绍 参见:https://blog.csdn.net/cs380637384/article/details/81217767
*/
//数据结构--文件名称
type Config struct {
	Name string //文件名
	Type string //文件类型
}

func init() {
	logs.Debug("1,init.....")
	InitConfig("config", "yaml")
}

func TestByViperInitConfig(t *testing.T) {
	logs.Debug("3,main......")
	router := gin.Default()
	router.GET("/test", testViperReadType)
	router.Run(":" + getServerPort())
}

//初始化配置
func InitConfig(cfg string, filetype string) {
	logs.Debug("2,initConfig.....")
	//实例化结构体
	c := Config{
		Name: cfg,
		Type: filetype,
	}
	//设置配置文件路径
	//viper.AddConfigPath("./resources") //此路径为相对于main.go的路径
	viper.AddConfigPath("./")
	//设置配置文件名称
	if c.Name != "" {
		//若指定了配置文件，则解析指定的配置文件
		viper.SetConfigName(c.Name)
	} else {
		//若未指定配置文件，则解析默认配置文件
		viper.SetConfigName("config")
	}
	//设置配置文件格式
	if c.Type != "" {
		viper.SetConfigType(c.Type)
	} else {
		viper.SetConfigType("yaml")
	}
	viper.AutomaticEnv()
	//viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("viper解析配置文件错误: %s \n", err))
	}
	//监听配置文件是否改变，用于热更新
	c.watchConfig()
}

//监听配置文件是否改变，用于热更新
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("resources file changed:%s\n", e.Name)
	})
}

//利用viper取string类型值
func getServerPort() string {
	serverPort := viper.GetString("server.port")
	logs.Info("serverPort:", serverPort)
	if serverPort == "" {
		return "8888"
	}
	return serverPort
}

//利用viper取数组值
func getArrayFromConfigYaml() {
	arr := viper.GetStringSlice("jvessel.peerserver")
	if arr != nil {
		for _, v := range arr {
			logs.Info("get array from config.yaml value:", v)
		}
	}
}

//测试viper读取的类型
func testViperReadType(ctx *gin.Context) {
	//1，viper从config.yaml读取数组
	getArrayFromConfigYaml()
}
