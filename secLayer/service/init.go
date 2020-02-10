package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ggkd/secLayer/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

// 初始化Waitgroup

var WG = sync.WaitGroup{}


// RedisBlackPool Redis连接池
var RedisLayer2ProxyPool *redis.Pool
var RedisProxy2LayerPool *redis.Pool

// 初始化RedisProxy2LayerPool
func InitRedisProxy2Layer() {
	RedisProxy2LayerPool = &redis.Pool{
		MaxIdle:     config.LayerConfig.RedisProxy2Layer.MaxIdle,
		MaxActive:   config.LayerConfig.RedisProxy2Layer.MaxActive,
		IdleTimeout: time.Duration(config.LayerConfig.RedisProxy2Layer.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf(config.LayerConfig.RedisProxy2Layer.Ip+":"+config.LayerConfig.RedisProxy2Layer.Port))
			return conn, err
		},
	}
	// 检测redis是否连接成功
	conn := RedisProxy2LayerPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		fmt.Println("RedisProxy2LayerPool err--------->", err)
		config.SugarLogger.Error("RedisProxy2LayerPool err--------->", err)
		return
	}
	fmt.Println("----------RedisProxy2LayerPool success----------")
	config.SugarLogger.Info("----------RedisProxy2LayerPool success----------")
}

// 初始化RedisLayer2ProxyPool
func InitRedisLayer2Proxy() {
	RedisLayer2ProxyPool = &redis.Pool{
		MaxIdle:     config.LayerConfig.RedisProxy2Layer.MaxIdle,
		MaxActive:   config.LayerConfig.RedisProxy2Layer.MaxActive,
		IdleTimeout: time.Duration(config.LayerConfig.RedisProxy2Layer.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf(config.LayerConfig.RedisProxy2Layer.Ip+":"+config.LayerConfig.RedisProxy2Layer.Port))
			return conn, err
		},
	}
	// 检测redis是否连接成功
	conn := RedisLayer2ProxyPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		fmt.Println("RedisLayer2ProxyPool err--------->", err)
		config.SugarLogger.Error("RedisLayer2ProxyPool err--------->", err)
		return
	}
	fmt.Println("----------RedisLayer2ProxyPool成功----------")
	config.SugarLogger.Info("----------RedisLayer2ProxyPool成功----------")
}


var ProxyEtcdClient *clientv3.Client // Etcd全局客户端
// 初始化etcd
func InitEtcd() {
	var err error
	endpoints := config.LayerConfig.Etcd.Ip + ":" + config.LayerConfig.Etcd.Port
	ProxyEtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoints},
		DialTimeout: time.Second * time.Duration(config.LayerConfig.Etcd.DialTimeout),
	})
	if err != nil {
		fmt.Println("初始化 ProxyEtcdClient err--------->", err)
		config.SugarLogger.Error("初始化 ProxyEtcdClient err--------->", err)
		return
	}
	fmt.Println("----------初始化ProxyEtcdClient成功----------")
	config.SugarLogger.Info("----------初始化ProxyEtcdClient成功----------")
}


// 监控etcd
func WatchEtcd() {
	fmt.Println("---------etcd watching-----------")
	config.SugarLogger.Debug("---------etcd watching-----------")
	key := fmt.Sprintf("%s/%s", config.LayerConfig.Etcd.SecKill_key, config.LayerConfig.Etcd.ProductKey)
	for {
		watchChan := ProxyEtcdClient.Watch(context.Background(), key)
		var getConfSuccess = true
		var secProductInfos []config.SecKillInfo
		for event := range watchChan {
			for _, ev := range event.Events {
				if ev.Type != mvccpb.DELETE {
					// 判断是否为删除事件
					err := json.Unmarshal(ev.Kv.Value, &secProductInfos)
					if err != nil {
						fmt.Println("unmarshal err : ", err)
						getConfSuccess = false
						continue
					}
				}
			}
			if getConfSuccess {
				UpdateSecProduct(secProductInfos)
			}
		}
	}
}


// 更新最新的商品信息
func UpdateSecProduct(productInfo []config.SecKillInfo)  {
	tmp := make(map[int]*config.SecKillInfo, 1000)
	for i, _ := range productInfo {
		tmp[productInfo[i].ProductId] = &productInfo[i]
		tmp[productInfo[i].ProductId].SecLimit = &SecLimit{}
	}
	config.LayerConfig.ProductRwLock.Lock()
	config.LayerConfig.SecKillProductMap = tmp
	config.LayerConfig.ProductRwLock.Unlock()
	fmt.Println("------------------------->", config.LayerConfig.SecKillProductMap)
	config.SugarLogger.Info("------------------------->", config.LayerConfig.SecKillProductMap)
}

// 获取秒杀商品的配置
func InitSecInfo()  {
	key := fmt.Sprintf("%s/%s", config.LayerConfig.Etcd.SecKill_key, config.LayerConfig.Etcd.ProductKey)
	resp, err := ProxyEtcdClient.Get(context.Background(), key)
	if err != nil {
		fmt.Println("get secInfo Key err--------->", err)
		config.SugarLogger.Error("get secInfo Key err--------->", err)
		return
	}
	var ProductInfos []config.SecKillInfo
	for _, kv := range resp.Kvs {
		//fmt.Printf("get secInfo Key[%v], Value[%v]\n", string(kv.Key), string(kv.Value))
		err := json.Unmarshal(kv.Value, &ProductInfos)
		if err != nil {
			fmt.Println("unmarshal err-------->", err)
			config.SugarLogger.Errorf("unmarshal err-------->", err)
			return
		}
		config.SugarLogger.Debugf("get secInfo Key[%v], Value[%v]", string(kv.Key), string(kv.Value))
	}
	UpdateSecProduct(ProductInfos)
}


// 全局初始化
func InitProxy() {
 	fmt.Println("===========初始化Proxy Log、Redis、Etcd===========")
	config.InitLog()
	InitRedisProxy2Layer()
	InitRedisLayer2Proxy()
	InitEtcd()
 	InitSecInfo()
	go WatchEtcd()
}
