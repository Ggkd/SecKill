package service

import "sync"

type UserHistory struct {
	History map[int]int
	Lock sync.RWMutex
}

func (u *UserHistory) Get(productId int) int {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	count, _ := u.History[productId]
	return count
}

func (u *UserHistory) Add(productId, count int) {
	u.Lock.Lock()
	defer u.Lock.Unlock()
	currentCount, ok := u.History[productId]
	if !ok {
		currentCount = count
	}else {
		currentCount += count
	}
	u.History[productId] = currentCount
}