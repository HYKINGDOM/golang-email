package sse

import (
    "sync"
)

type Event struct {
    TaskID uint   `json:"task_id"`
    Type   string `json:"type"`
    Data   any    `json:"data"`
}

type Broker struct {
    mu      sync.RWMutex
    clients map[uint][]chan Event
}

func NewBroker() *Broker { return &Broker{clients: make(map[uint][]chan Event)} }

func (b *Broker) Subscribe(taskID uint) chan Event {
    ch := make(chan Event, 16)
    b.mu.Lock()
    b.clients[taskID] = append(b.clients[taskID], ch)
    b.mu.Unlock()
    return ch
}

func (b *Broker) Unsubscribe(taskID uint, ch chan Event) {
    b.mu.Lock()
    defer b.mu.Unlock()
    arr := b.clients[taskID]
    for i, c := range arr {
        if c == ch {
            b.clients[taskID] = append(arr[:i], arr[i+1:]...)
            close(c)
            break
        }
    }
}

func (b *Broker) Publish(taskID uint, evt Event) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    for _, ch := range b.clients[taskID] {
        select { case ch <- evt: default: }
    }
}