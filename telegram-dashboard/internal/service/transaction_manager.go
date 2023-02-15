package service

import (
	"fmt"
	"sync"
	"time"
	"yt-solutions-telegram-dashboard/internal/core"
)

type TransactionManager interface {
	Set(tx *core.Transaction)
	Get(id string) (*core.Transaction, bool)
	Delete(id ...string)
}

type transactionManager struct {
	cleanInterval int
	mu            sync.RWMutex
	storage       map[string]*core.Transaction
}

func NewTransactionManager(interval int) TransactionManager {
	t := &transactionManager{
		cleanInterval: interval,
		mu:            sync.RWMutex{},
		storage:       make(map[string]*core.Transaction),
	}

	go t.StartGC()

	return t
}

func (t *transactionManager) Set(tx *core.Transaction) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.storage[tx.Id] = tx

	fmt.Println(t.storage)
}

func (t *transactionManager) Get(id string) (*core.Transaction, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	tx, ok := t.storage[id]
	return tx, ok
}

func (t *transactionManager) Delete(id ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, tx := range id {
		delete(t.storage, tx)
	}
}

func (t *transactionManager) StartGC() {
	for {
		<-time.After(time.Duration(t.cleanInterval) * time.Second)
		expired := t.expiredTransactions()
		t.Delete(expired...)
	}
}

func (t *transactionManager) expiredTransactions() []string {

	t.mu.RLock()
	defer t.mu.RUnlock()

	var expired []string
	for _, tx := range t.storage {
		if tx.IsExpired() {
			expired = append(expired, tx.Id)
		}
	}

	return expired
}
