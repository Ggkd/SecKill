package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func PutEtcdProduct(client *clientv3.Client)  {
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	type SecKillInfo struct {
		ProductId int
		StartTime int
		EndTime   int
		Count     int
		Status    int
	}
	p1 := SecKillInfo{
		ProductId: 11,
		StartTime: 1580711459,
		EndTime:   1580712878,
		Count:     100,
		Status:    1,
	}
	p2 := SecKillInfo{
		ProductId: 111,
		StartTime: 1580711459,
		EndTime:   1580712878,
		Count:     50,
		Status:    1,
	}
	var ProductInfos []SecKillInfo
	ProductInfos = append(ProductInfos, p1, p2)
	key := "secKill/product"
	data, err := json.Marshal(ProductInfos)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	_, err = client.Put(ctx, key, string(data))
	if err != nil {
		fmt.Println(err)
		return
	}
}


func main()  {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:time.Second*3,
	})
	if err != nil {
		fmt.Println("new client err :", err)
		return
	}
	fmt.Println("conn etcd success")
	PutEtcdProduct(client)
}