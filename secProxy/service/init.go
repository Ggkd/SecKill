package service

import (
	"fmt"
	"github.com/Ggkd/secProxy/conf"
	"github.com/garyburd/redigo/redis"
	"time"
)

// 获取黑名单
func GetRedisBlackList()  {
	conn := conf.RedisBlackPool.Get()
	// 获取黑名单的id list
	reply, err := conn.Do("hgetall", "idblacklist")
	idlist, err := redis.Strings(reply, err)
	if err != nil {
		conf.SugarLogger.Error("hgetall user idblacklist err-------->", err)
		fmt.Println("hget user idblacklist err-------->", err)
		return
	}
	for _, id := range idlist {
		conf.Config.UserIdBlackList[id] = true
	}
	// 获取黑名单的ip list
	reply, err = conn.Do("hgetall", "ipblacklist")
	iplist, err := redis.Strings(reply, err)
	if err != nil {
		conf.SugarLogger.Error("hgetall user ipblacklist err-------->", err)
		fmt.Println("hget user ipblacklist err-------->", err)
		return
	}
	for _, id := range iplist {
		conf.Config.UserIpBlackList[id] = true
	}
	conn.Close()
	go syncBlackIdList()
	go syncBlackIpList()
}

//同步id黑名单
func syncBlackIdList()  {
	var idList []string
	var lastTime = time.Now().Unix()
	for {
		conn := conf.RedisBlackPool.Get()
		// 获取黑名单的id list
		reply, err := conn.Do("BLPOP", "idblacklist", 1)
		id, err := redis.String(reply, err)
		if err != nil {
			conf.SugarLogger.Error("BLPOP user idblacklist err-------->", err)
			fmt.Println("BLPOP user idblacklist err-------->", err)
			continue
		}
		idList = append(idList, id)
		nowTime := time.Now().Unix()
		if len(idList) > 100 || nowTime - lastTime > 5 {
			conf.Config.BlackRwLock.Lock()
			for _, id := range idList {
				conf.Config.UserIdBlackList[id] = true
			}
			conf.Config.BlackRwLock.Unlock()
			lastTime = nowTime
		}
	}
}
//同步ip黑名单
func syncBlackIpList()  {
	var ipList []string
	var lastTime = time.Now().Unix()
	for {
		conn := conf.RedisBlackPool.Get()
		// 获取黑名单的ip list
		reply, err := conn.Do("BLPOP", "ipblacklist", 1)
		ip, err := redis.String(reply, err)
		if err != nil {
			conf.SugarLogger.Error("BLPOP user ipblacklist err-------->", err)
			fmt.Println("BLPOP user ipblacklist err-------->", err)
			continue
		}
		ipList = append(ipList, ip)
		nowTime := time.Now().Unix()
		if len(ipList) > 100 || nowTime - lastTime > 5 {
			conf.Config.BlackRwLock.Lock()
			for _, id := range ipList {
				conf.Config.UserIpBlackList[id] = true
			}
			conf.Config.BlackRwLock.Unlock()
			lastTime = nowTime
		}
	}
}

// 初始化redis处理队列
func InitRedisProcess()  {
	for i:=0; i< conf.Config.RedisProxy2Layer.WriteGoroutineNum; i++ {
		go WriteHandle()
	}

	for i:=0; i< conf.Config.RedisProxy2Layer.ReadGoroutineNum; i++ {
		go ReadHandle()
	}
}


func init()  {
	conf.Config.RedisProxy2Layer.ReqChan = make(chan *conf.ReqSecKill, conf.Config.RedisProxy2Layer.ReqChanSize)
	GetRedisBlackList()
	syncBlackIdList()
	syncBlackIpList()
	InitRedisProcess()
}