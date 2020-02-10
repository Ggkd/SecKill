package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ggkd/secLayer/config"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

func Run()  {
	RunProcess()
}

// 进程调度
func RunProcess()  {
	for i:=0; i< config.LayerConfig.Service.ReadGoroutineNum; i++ {
		WG.Add(1)
		go HandleRead()
	}

	for i:=0; i< config.LayerConfig.Service.WriteGoroutineNum; i++ {
		WG.Add(1)
		go HandleWrite()
	}

	for i:=0; i< config.LayerConfig.Service.HandleGoroutineNum; i++ {
		WG.Add(1)
		go HandleUser()
	}

	fmt.Println("all process goroutine start")
	config.SugarLogger.Info("all process goroutine start")
	WG.Wait()
	fmt.Println("all process goroutine end")
	config.SugarLogger.Info("all process goroutine end")
}


// 读处理
func HandleRead()  {
	for {
		// 从redis pool 获取一个连接
		conn := RedisProxy2LayerPool.Get()
		for {
			reply, err := conn.Do("blpop", "queuelist", 0)
			data, err := redis.String(reply, err)
			if err != nil {
				fmt.Println("blpop err--------->", err)
				config.SugarLogger.Error("blpop err--------->", err)
				break
			}

			config.SugarLogger.Debugf("blpop value [%v]\n", data)

			var req config.ReqSecKill
			err = json.Unmarshal([]byte(data), &req)
			if err != nil {
				fmt.Println("unmarshal req data err-------->", err)
				config.SugarLogger.Error("unmarshal req data err-------->",  err)
				continue
			}
			// 判断用户的请求是否超时
			nowTime := time.Now().Unix()
			if nowTime - req.AccessTime.Unix() > int64(config.LayerConfig.Service.MaxTimeOut) {
				config.SugarLogger.Warn("[%v] req timeout", req)
				continue
			}
			// 将请求发送到通道
			timer := time.NewTimer(time.Duration(config.LayerConfig.Service.ChanWaitTime) * time.Millisecond)
			select {
			case config.LayerConfig.ReadHandleChan <- &req:
			case <- timer.C:
				fmt.Println("send to readChan timeout")
				config.SugarLogger.Warn("send to readChan timeout")
			}

		}
	}
}

// 写处理
func HandleWrite()  {

}

//用户处理
func HandleUser()  {
	// 从通道取数据
	fmt.Println("start handle user req")
	config.SugarLogger.Debug("start handle user req")
	for req := range config.LayerConfig.ReadHandleChan {
		fmt.Printf("running req[%v]\n", req)
		config.SugarLogger.Debugf("running req[%v]\n", req)
		resp := handleSecKill(req)
		// 将响应发送写通道
		timer := time.NewTicker(time.Duration(config.LayerConfig.Service.ChanWaitTime) * time.Millisecond)
		select {
		case config.LayerConfig.WriteHandleChan <- resp:
		case <- timer.C:
			fmt.Println("send to writeChan timeout")
			config.SugarLogger.Warn("send to writeChan timeout")
		}
	}
}


// 具体的秒杀处理逻辑
func handleSecKill(req *config.ReqSecKill) *config.RespSecKill {
	// 判断商品是否存在
	var resp =  &config.RespSecKill{}
	product, ok := config.LayerConfig.SecKillProductMap[req.ProductId]
	if !ok {
		resp.Code = 0
		resp.Msg = "product not found"
		return resp
	}
	// 判断商品的状态
	if product.Status == 0 {
		resp.Code = 0
		resp.Msg = "product sale out"
		return resp
	}
	// 校验商品每秒的售出数据是否大于限制的数量
	nowTime := time.Now().Unix()
	if product.SecLimit.Check(nowTime) >= product.MaxSecNum {
		resp.Code = 0
		resp.Msg = "刷新重试"
		return resp
	}
	// 判断用户购买该商品的数量是否超出限制
	config.LayerConfig.UserBuyHistoryLock.Lock()
	productHistory, ok := config.LayerConfig.UserBuyHistory[req.UserId]
	if !ok {
		productHistory = &UserHistory{History:make(map[int]int, 16)}
		config.LayerConfig.UserBuyHistory[req.UserId] = productHistory
	}
	count := productHistory.Get(product.ProductId)
	config.LayerConfig.UserBuyHistoryLock.Unlock()
	if count >= product.UserBuyLimit {
		resp.Code = 0
		resp.Msg = "已超出购买数量"
		return resp
	}
	// 判断商品是否售空
	if product.Count == 0{
		resp.Code = 0
		resp.Msg = "商品已售完"
		return resp
	}
	// 判断用户买到的概率
	rate := rand.Float64()
	if rate < product.BuyRate {
		resp.Code = 0
		resp.Msg = "刷新重试"
		return resp
	}
	// 用户可购买
	productHistory.Add(product.ProductId, 1)
	config.LayerConfig.ProductRwLock.Lock()
	product.Count -= 1
	config.LayerConfig.ProductRwLock.Unlock()
	resp.Code = 1
	resp.ProductId = product.ProductId
	resp.UserId = req.UserId
	resp.Msg = "购买成功"
	return resp
}