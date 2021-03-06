package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/Ggkd/secProxy/conf"
	"time"
)

// 当前商品状态信息
type currentProductInfo struct {
	ProductId int
	Start     bool
	End       bool
	Status    string
}

//获取查询的商品
func SecInfoById(productId int) interface{} {
	var productList []currentProductInfo
	currentProduct := checkProduct(productId)
	productList = append(productList, currentProduct)
	return productList
}

//获取所有的商品
func SecInfoList() interface{} {
	var productList []currentProductInfo
	conf.Config.ProductRwLock.RLock()
	for id, _ := range conf.Config.SecKillProductMap {
		currentProduct := checkProduct(id)
		productList = append(productList, currentProduct)
	}
	conf.Config.ProductRwLock.RUnlock()
	return productList
}

//获取商品的状态
func checkProduct(productId int) currentProductInfo {
	currentProduct := currentProductInfo{}
	conf.Config.ProductRwLock.RLock()
	v, ok := conf.Config.SecKillProductMap[productId]
	conf.Config.ProductRwLock.RUnlock()
	if !ok {
		currentProduct.ProductId = productId
		currentProduct.Status = "product isn't exist"
		return currentProduct
	}
	now := time.Now().Unix()
	start := false
	end := false
	status := "secKill is starting"
	if now < v.StartTime {
		status = "secKill not start"
	}
	if now > v.StartTime {
		start = true
	}
	if now > v.EndTime {
		start = false
		end = true
		status = "secKill is ended"
	}
	if v.Status != 1 {
		status = "product is sale out"
	}
	currentProduct.ProductId = productId
	currentProduct.Start = start
	currentProduct.End = end
	currentProduct.Status = status
	fmt.Println(currentProduct)
	return currentProduct
}



func SecKill(req *conf.ReqSecKill) (map[string]interface{}, int, error) {
	data := make(map[string]interface{})
	err := checkUser(req)
	if err != nil {
		conf.SugarLogger.Error(err)
		return data, 0, err
	}
	// 校验用户是否非法访问
	err = Antispam(req)
	if err != nil {
		return data, 0 ,err
	}
	return data, 1, err
}

// 校验用户
func checkUser(req *conf.ReqSecKill) error {
	var err error
	// 校验用户是否登录
	userData := fmt.Sprintf("%s:%s", req.UserId, conf.Config.UserControl.Secret)
	userSign := fmt.Sprintf("%x", md5.Sum([]byte(userData)))
	if req.UserAuth != userSign {
		err = errors.New("invalid user cookie")
		return err
	}
	// 校验用户是否正常访问
	found := false
	for _, refer :=  range conf.Config.UserControl.ReferList {
		if req.UserRefer == refer {
			found = true
			break
		}
	}
	if !found {
		err = errors.New("invalid request")
		return err
	}
	conf.Config.RedisProxy2Layer.ReqChan <- req
	return err
}