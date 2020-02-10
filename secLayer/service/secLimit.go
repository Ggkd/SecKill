package service


// 商品每秒的售出数量
type SecLimit struct {
	count int
	currentTime int64
}

// 计数
func (ur *SecLimit) Count(accessTime int64) int {
	if ur.currentTime != accessTime {
		ur.count = 1
		ur.currentTime = accessTime
		return ur.count
	}
	ur.count ++
	return ur.count
}

// 获取当前访问时的数量
func (ur *SecLimit) Check(nowTime int64) int {
	if ur.currentTime != nowTime {
		return 0
	}
	return  ur.count
}