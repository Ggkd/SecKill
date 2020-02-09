package service

import (
	"errors"
	"github.com/Ggkd/secProxy/conf"
	"sync"
)

// 记录所有用户的访问记录
type UserReqMgr struct {
	UserReq map[string]*UserRecord
	UserIp	map[string]*UserRecord
	lock sync.Mutex
}

// 记录用户的访问频率
type UserRecord struct {
	count int
	currentTime int64
}

var userReqMgr = &UserReqMgr{UserReq:make(map[string]*UserRecord, 1000), UserIp:make(map[string]*UserRecord, 1000)}

// 校验用户的行为
func Antispam(req *conf.ReqSecKill) error {
	var err error
	userReqMgr.lock.Lock()
	// 校验用户的账号是否正常访问
	userIdRecord, ok := userReqMgr.UserReq[req.UserId]
	if !ok {
		userIdRecord = &UserRecord{}
		userReqMgr.UserReq[req.UserId] = userIdRecord
	}
	userIdCount := userIdRecord.Count(req.AccessTime.Unix())
	// 校验用户的ip是否正常访问
	userIpRecord, ok := userReqMgr.UserReq[req.UserAddr]
	if !ok {
		userIpRecord = &UserRecord{}
		userReqMgr.UserReq[req.UserAddr] = userIpRecord
	}
	userReqMgr.lock.Unlock()
	userIpCount := userIpRecord.Count(req.AccessTime.Unix())
	if userIdCount > conf.Config.UserControl.ReqLimit || userIpCount > conf.Config.UserControl.IpLimit {
		err = errors.New("service is busy")
		return err
	}
	return err
}

// 对用户的访问计数
func (ur *UserRecord) Count(accessTime int64) int {
	if ur.currentTime != accessTime {
		ur.count = 1
		ur.currentTime = accessTime
		return ur.count
	}
	ur.count ++
	return ur.count
}