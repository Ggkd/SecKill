package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/garyburd/redigo/redis"
	"time"
)

// RedisPool Redis连接池
var RedisPool *redis.Pool

// 初始化redis
func InitRedis()  {
	RedisPool = &redis.Pool{
		MaxIdle:     Config.Redis.MaxIdle,
		MaxActive:   Config.Redis.MaxActive,
		IdleTimeout: time.Duration(Config.Redis.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf(Config.Redis.Ip + ":" + Config.Redis.Port))
			return conn, err
		},
	}
	// 检测redis是否连接成功
	conn := RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		fmt.Println("初始化redis err--------->", err)
		SugarLogger.Error("初始化redis err--------->", err)
		return
	}
	fmt.Println("----------初始化redis成功----------")
	SugarLogger.Info("----------初始化redis成功----------")
}


var EtcdClient *clientv3.Client		// Etcd全局客户端
// 初始化etcd
func InitEtcd()  {
	var err error
	endpoints := Config.Etcd.Ip + ":" + Config.Etcd.Port
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:[]string{endpoints},
		DialTimeout:time.Second*time.Duration(Config.Etcd.DialTimeout),
	})
	if err != nil {
		fmt.Println("初始化etcd err--------->", err)
		SugarLogger.Error("初始化etcd err--------->", err)
		return
	}
	fmt.Println("----------初始化etcd成功----------")
	SugarLogger.Info("----------初始化etcd成功----------")
}


// 获取秒杀商品的配置
func InitSecInfo()  {
	key := fmt.Sprintf("%s/%s", Config.Etcd.SecKill_key, Config.Etcd.ProductKey)
	resp, err := EtcdClient.Get(context.Background(), key)
	if err != nil {
		fmt.Println("get secInfo Key err--------->", err)
		SugarLogger.Error("get secInfo Key err--------->", err)
		return
	}
	var ProductInfos []SecKillInfo
	for _, kv := range resp.Kvs {
		//fmt.Printf("get secInfo Key[%v], Value[%v]\n", string(kv.Key), string(kv.Value))
		err := json.Unmarshal(kv.Value, &ProductInfos)
		if err != nil {
			fmt.Println("unmarshal err-------->", err)
			SugarLogger.Errorf("unmarshal err-------->", err)
			return
		}
		SugarLogger.Debugf("get secInfo Key[%v], Value[%v]", string(kv.Key), string(kv.Value))
	}
	UpdateSecProduct(ProductInfos)
}


// 监控etcd
func WatchEtcd()  {
	fmt.Println("---------etcd watching-----------")
	SugarLogger.Debug("---------etcd watching-----------")
	key := fmt.Sprintf("%s/%s", Config.Etcd.SecKill_key, Config.Etcd.ProductKey)
	for {
		watchChan := EtcdClient.Watch(context.Background(), key)
		var getConfSuccess = true
		var secProductInfos []SecKillInfo
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
func UpdateSecProduct(productInfo []SecKillInfo)  {
	tmp := make(map[int]*SecKillInfo, 1000)
	for i, _ := range productInfo {
		tmp[productInfo[i].ProductId] = &productInfo[i]
	}
	Config.RwLock.Lock()
	Config.SecKillProductMap = tmp
	Config.RwLock.Unlock()
	fmt.Println("------------------------->", Config.SecKillProductMap)
	SugarLogger.Info("------------------------->", Config.SecKillProductMap)
}


// 全局初始化
func init()  {
	Config.SecKillProductMap = make(map[int]*SecKillInfo, 1000)
	fmt.Println("===========初始化Log、Redis、Etcd===========")
	InitLog()
	InitRedis()
	InitEtcd()
	InitSecInfo()
	go WatchEtcd()
}