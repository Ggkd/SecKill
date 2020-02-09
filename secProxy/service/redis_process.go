package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ggkd/secProxy/conf"
)

// 写处理
func WriteHandle()  {
	for {
		req := <- conf.Config.RedisProxy2Layer.ReqChan
		conn := conf.RedisProxy2LayerPool.Get()
		data, err := json.Marshal(req)
		if err != nil {
			conf.SugarLogger.Error("marshal req err------>", err)
			conn.Close()
			continue
		}
		_, err = conn.Do("LPUSH", "sec_queue", data)
		if err != nil {
			fmt.Println("Lpush req err ----------->", err)
			conn.Close()
			continue
		}
		conn.Close()
	}
}

// 读处理
func ReadHandle()  {

}