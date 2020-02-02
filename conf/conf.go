package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

// 总配置
type Conf struct {
	Host  `ini:"host"`
	Redis `ini:"redis"`
	Etcd  `ini:"etcd"`
	Log   `ini:"log"`
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
}

// log配置
type Log struct {
	Path  string `ini:"path"`
	Level string `ini:"level"`
}

// 秒杀商品配置
type SecKillInfo struct {
	ProductId 	int
	StartTime	time.Time
	EndTime		time.Time
	Count		int
	Status 		int
}


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
