package service

import (
	"errors"
	"fmt"
	"github.com/Ggkd/conf"
)

func SecInfoService(productId int) (conf.SecKillInfo, error) {
	var err error
	fmt.Println(productId)
	conf.Config.RwLock.RLock()
	v, ok := conf.Config.SecKillProductMap[productId]
	fmt.Println(v)
	conf.Config.RwLock.RUnlock()
	if !ok {
		err = errors.New("product isn't exist")
		return conf.SecKillInfo{}, err
	}
	return *v, err
}