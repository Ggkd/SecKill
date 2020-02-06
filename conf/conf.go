package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

// 总配置
type Conf struct {
	Host              `ini:"host"`
	Redis             `ini:"redis"`
	Etcd              `ini:"etcd"`
	Log               `ini:"log"`
	SecKillProductMap map[int]*SecKillInfo
	RwLock            sync.RWMutex
	UserControl       `ini:"user_control"`
	UserIdBlackList   map[string]bool
	UserIpBlackList   map[string]bool
}

// 主机配置
type Host struct {
	Ip   string `ini:"ip"`
	Port string `ini:"port"`
}

// redis配置
type Redis struct {
	Ip          string `ini:"ip"`
	Port        string `ini:"port"`
	MaxIdle     int    `ini:"MaxIdle"`
	MaxActive   int    `ini:"MaxActive"`
	IdleTimeout int    `ini:"IdleTimeout"`
}

// etcd配置
type Etcd struct {
	Ip          string `ini:"ip"`
	Port        string `ini:"port"`
	DialTimeout int    `ini:"DialTimeout"`
	SecKill_key string `ini:"secKill_key"`
	ProductKey  string `ini:"product_key"`
}

// log配置
type Log struct {
	Path  string `ini:"path"`
	Level string `ini:"level"`
}

// 秒杀商品配置
type SecKillInfo struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Count     int
	Status    int
}

// 用户请求配置
type ReqSecKill struct {
	ProductId  int
	UserId     string
	SecTime    int
	Source     string
	Nance      string
	AuthCode   string
	UserAuth   string
	AccessTime time.Time
	UserAddr   string
	UserRefer  string
}

type UserControl struct {
	Secret    string   `int:"secret"`
	ReqLimit  int      `ini:"req_limit"`
	IpLimit   int      `ini:"ip_limit"`
	ReferList []string `ini:"refer_list"`
}

//返回响应
type Result struct {
	Code int
	Msg  string
	Data interface{}
}

// 全局配置对象
var Config = new(Conf)

//从配置文件加载所有配置
func init() {
	fmt.Println("===========从配置文件加载所有配置===========")
	err := ini.MapTo(Config, "../conf/conf.ini")
	if err != nil {
		fmt.Println("加载配置err--------->", err)
		return
	}
	fmt.Println("----------加载配置成功----------")
}
