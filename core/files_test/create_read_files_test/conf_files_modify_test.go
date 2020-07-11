package modify_files

import (
	"errors"
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/logs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

/*
 示例：修改.conf文件 add by syf 2020.5.26
 说明：应用库github.com/Unknwon/goconfig
*/
type LedgerInit struct {
	NetworkName      string         `json:"network_name"`      //网络名称
	LedgerSeed       string         `json:"ledger_seed"`       //账本种子(创建账本需要修改该值，初始化账本默认即可)
	LedgerName       string         `json:"ledger_name"`       //账本名字
	UserName         string         `json:"user_name"`         //当前参与者
	CreateTime       string         `json:"create_time"`       //创建时间
	ConsPartiCount   int            `json:"cons_parti_count"`  //共识参与方数量
	SecurityRoles    string         `json:"security_roles"`    //角色
	LedgerPrivileges string         `json:"ledger_privileges"` //账本权限
	TxPrivileges     string         `json:"tx_privileges"`     //交易权限
	LedgerNum        int            `json:"ledgerNum"`         //账本数量
	Participants     []Participants `json:"participants"`      //参与方
}

//参与方
type Participants struct {
	ParticipantsName        string `json:"participants_name"`         //参与方名字
	ParticipantsType        string `json:"participants_type"`         //参与方类型(leader:发起者；inviter:受邀者)
	PubKeyPath              string `json:"pubkey_path"`               //公钥路径
	PubKey                  string `json:"pubkey"`                    //公钥
	Roles                   string `json:"roles"`                     //角色
	RolesPolicy             string `json:"roles_policy"`              //角色权限策略
	ConsensusHost           string `json:"consensus_host"`            //账本初始化共识主机IP
	ConsensusPort           string `json:"consensus_port"`            //账本初始化共识端口(8900)
	ConsensusSecure         string `json:"consensus_secure"`          //账本初始化共识服务是否开启安全连接
	InitializerHost         string `json:"initializer_host"`          //账本初始化主机IP
	InitializerPort         string `json:"initializer_port"`          //账本初始化端口(8800)
	InitializerSecure       string `json:"initializer_secure"`        //账本初始化是否开启安全连接
	ConsensusBftsmartHost   string `json:"consensus_bftsmart_host"`   //bftsmart共识主机
	ConsensusBftsmartPort   string `json:"consensus_bftsmart_port"`   //bftsmart共识端口(16000)
	ConsensusBftsmartSecure string `json:"consensus_bftsmart_secure"` //bftsmart共识服务是否开启安全连接
	PeerIp                  string `json:"peer_ip"`                   //peer服务ip
	PeerPort                string `json:"peer_port"`                 //peer服务端口(7080)
	IsGateway               bool   `json:"is_gateway"`                //是否安装Gateway(true:是；false:否)
}

//示例：修改/files/ledger.init配置文件
func TestModifyLedgerInit(t *testing.T) {
	params := LedgerInit{
		NetworkName:      "syfnet",
		LedgerSeed:       "syfseed",
		LedgerName:       "syfledger",
		UserName:         "syf",
		CreateTime:       "2020-5-26 12:00:00",
		ConsPartiCount:   2,
		SecurityRoles:    "1",
		LedgerPrivileges: "1",
		TxPrivileges:     "1",
		LedgerNum:        1,
		Participants:     nil,
	}
	userArr := []Participants{}
	user := Participants{
		ParticipantsName:        "zs",
		ParticipantsType:        "1",
		PubKeyPath:              "",
		PubKey:                  "123456",
		Roles:                   "1",
		RolesPolicy:             "1",
		ConsensusHost:           "127.0.0.1",
		ConsensusPort:           "8080",
		ConsensusSecure:         "false",
		InitializerHost:         "127.0.0.1",
		InitializerPort:         "8081",
		InitializerSecure:       "false",
		ConsensusBftsmartHost:   "127.0.0.1",
		ConsensusBftsmartPort:   "8082",
		ConsensusBftsmartSecure: "false",
		PeerIp:                  "127.0.0.1",
		PeerPort:                "7080",
		IsGateway:               false,
	}
	userArr = append(userArr, user)
	params.Participants = userArr
	filename_ledgerinit := filepath.Join("/tmp", "ledger.init")
	filename_bftsmart := filepath.Join("/tmp", "bftsmart.config")
	if err := ModifyLedgerInit(params, filename_ledgerinit, filename_bftsmart); err != nil {
		logs.Error("modify files fail,", err)
		return
	}
}

func ModifyLedgerInit(params LedgerInit, filesname_ledgerInit, filesname_bftsmart string) error {
	logs.Info("modify ledger.init config...filesname=", filesname_ledgerInit)
	//加载ledger.init
	cfg_ledgerInit, err := goconfig.LoadConfigFile(filesname_ledgerInit)
	if err != nil {
		logs.Error("load ledger.init config file fail,err:", err.Error())
		return err
	}
	if cfg_ledgerInit == nil {
		logs.Error("load ledger.init config is null!")
		return errors.New("加载ledger.init配置文件为空!")
	}
	//加载bftsmart
	cfg_bftsmart, err := goconfig.LoadConfigFile(filesname_bftsmart)
	if err != nil {
		logs.Error("load bftsmart.config fail,err:", err.Error())
		return err
	}
	if cfg_bftsmart == nil {
		logs.Error("load bftsmart.config is null!")
		return errors.New("加载bftsmart.config配置文件为空!")
	}
	//设置ledger.init参数值
	if params.LedgerSeed != "" {
		cfg_ledgerInit.SetValue("", "ledger.seed", params.LedgerSeed)
	}
	cfg_ledgerInit.SetValue("", "ledger.name", params.LedgerName)
	cfg_ledgerInit.SetValue("", "created-time", params.CreateTime)
	cfg_ledgerInit.SetValue("", "cons_parti.count", strconv.Itoa(len(params.Participants)))
	//bftsmart参数值
	var initialView string
	//参与者参数设置
	for k, v := range params.Participants {
		//ledger.init参数值
		cfg_ledgerInit.SetKeyComments("", "cons_parti."+strconv.Itoa(k)+".name", "# 第"+strconv.Itoa(k)+"个账本参与方")
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".name", v.ParticipantsName)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".pubkey", v.PubKey)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".consensus.host", v.ConsensusHost)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".consensus.port", v.ConsensusPort)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".consensus.secure", v.ConsensusSecure)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".initializer.host", v.InitializerHost)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".initializer.port", v.InitializerPort)
		cfg_ledgerInit.SetValue("", "cons_parti."+strconv.Itoa(k)+".initializer.secure", v.InitializerSecure)

		//bftsmart参数值
		cfg_bftsmart.SetKeyComments("", "system.server."+strconv.Itoa(k)+".network.host", "# 第"+strconv.Itoa(k)+"个共识参与方")
		cfg_bftsmart.SetValue("", "system.server."+strconv.Itoa(k)+".network.host", v.ConsensusBftsmartHost)
		cfg_bftsmart.SetValue("", "system.server."+strconv.Itoa(k)+".network.port", v.ConsensusBftsmartPort)
		cfg_bftsmart.SetValue("", "system.server."+strconv.Itoa(k)+".network.secure", v.ConsensusBftsmartSecure)
		initialView = initialView + strconv.Itoa(k) + ","
	}
	cfg_bftsmart.SetValue("", "system.servers.num", strconv.Itoa(len(params.Participants)))
	cfg_bftsmart.SetValue("", "system.initial.view", strings.TrimRight(initialView, ","))
	//更改ledger.init参数值
	if err = goconfig.SaveConfigFile(cfg_ledgerInit, filesname_ledgerInit); err != nil {
		logs.Error("modify ledger.init config fail,err:", err.Error())
		return err
	}
	//更改bftsmart参数值
	if err = goconfig.SaveConfigFile(cfg_bftsmart, filesname_bftsmart); err != nil {
		logs.Error("modify bftsmart.config fail,err:", err.Error())
		return err
	}
	return nil
}
