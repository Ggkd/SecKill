package config

import (
	"fmt"
	"github.com/Ggkd/secLayer/service"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

// 总配置
type Conf struct {
	Host               `ini:"host"`
	Etcd               `ini:"etcd"`
	Log                `ini:"log"`
	ProductRwLock      sync.RWMutex
	UserIdBlackList    map[string]bool
	UserIpBlackList    map[string]bool
	BlackRwLock        sync.RWMutex
	RedisProxy2Layer   `ini:"redis_proxy2layer"`
	SecKillProductMap  map[int]*SecKillInfo
	Service            `ini:"service"`
	ReadHandleChan     chan *ReqSecKill
	WriteHandleChan    chan *RespSecKill
	UserBuyHistory     map[string]*service.UserHistory
	UserBuyHistoryLock sync.RWMutex
}

// 主机配置
type Host struct {
	Ip   string `ini:"ip"`
	Port string `ini:"port"`
}

// redis接口层——>逻辑层
type RedisProxy2Layer struct {
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

// 服务处理配置
type Service struct {
	WriteGoroutineNum   int `ini:"writeGoroutineNum"`
	ReadGoroutineNum    int `ini:"readGoroutineNum"`
	HandleGoroutineNum  int `ini:"handleGoroutineNum"`
	ReadHandleChanSize  int `ini:"readHandleChanSize"`
	WriteHandleChanSize int `ini:"writeHandleChanSize"`
	MaxTimeOut          int `ini:"MaxTimeOut"`
	ChanWaitTime        int `ini:"ChanWaitTime"`
}

// 秒杀商品配置
type SecKillInfo struct {
	ProductId    int
	StartTime    int64
	EndTime      int64
	Count        int
	Status       int
	MaxSecNum    int // 每秒最多售出的数量
	SecLimit     *service.SecLimit
	UserBuyLimit int     // 用户最多购买的数量
	BuyRate      float64 //用户买到的概率
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

// 返回响应配置
type RespSecKill struct {
	ProductId int
	UserId    string
	Token     string
	Code      int
	Msg       string
}

// 全局配置对象
var LayerConfig = new(Conf)

//从配置文件加载所有配置
func init() {
	fmt.Println("===========从配置文件加载所有Layer配置===========")
	err := ini.MapTo(LayerConfig, "../config/layer_conf.ini")
	if err != nil {
		fmt.Println("加载Layer配置err--------->", err)
		return
	}
	LayerConfig.ReadHandleChan = make(chan *ReqSecKill, LayerConfig.Service.ReadHandleChanSize)
	LayerConfig.WriteHandleChan = make(chan *RespSecKill, LayerConfig.Service.WriteHandleChanSize)
	LayerConfig.UserBuyHistory = make(map[int]*SecKillInfo, 10000)
	fmt.Println("----------加载Layer配置成功----------")
}
