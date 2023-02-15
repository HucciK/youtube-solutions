package service

import (
	"sync"
	"yt-solutions-telegram-dashboard/internal/core"
)

type StateMachine interface {
	Set(user *core.User)
	Get(userId int) (*core.User, bool)
}

type stateMachine struct {
	mu      sync.RWMutex
	storage map[int]*core.User
}

func NewStateMachine() StateMachine {
	return &stateMachine{
		mu:      sync.RWMutex{},
		storage: make(map[int]*core.User),
	}
}

func (sm stateMachine) Set(user *core.User) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.storage[user.ID] = user
}

func (sm stateMachine) Get(userId int) (*core.User, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	user, ok := sm.storage[userId]

	return user, ok
}
