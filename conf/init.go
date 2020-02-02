package conf

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
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
		sugarLogger.Error("初始化redis err--------->", err)
		return
	}
	fmt.Println("----------初始化redis成功----------")
	sugarLogger.Info("----------初始化redis成功----------")
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
		sugarLogger.Error("初始化etcd err--------->", err)
		return
	}
	fmt.Println("----------初始化etcd成功----------")
	sugarLogger.Info("----------初始化etcd成功----------")
}

func InitSecInfo()  {
	resp, err := EtcdClient.Get(context.Background(), Config.Etcd.SecKill_key)
	if err != nil {
		fmt.Println("get secInfo Key err--------->", err)
		sugarLogger.Error("get secInfo Key err--------->", err)
		return
	}
	for k, v := range resp.Kvs {
		sugarLogger.Debugf("get secInfo Key[%v], Value[%v]", k, v)
	}

}

func init()  {
	fmt.Println("===========初始化Log、Redis、Etcd===========")
	InitLog()
	InitRedis()
	InitEtcd()
}