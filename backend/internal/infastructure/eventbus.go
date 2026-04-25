package infastructure

import(
	"sync"
)

type EventBus struct {
    mu    sync.RWMutex
    subscribers map[string][]chan any
}

func NewEventBus() *EventBus {
    return &EventBus{
        subscribers: make(map[string][]chan any),
    }
}

func (eb *EventBus) Publish(topic string, event any) {
    eb.mu.RLock()
    defer eb.mu.RUnlock()

    subscribers := append([]chan any{}, eb.subscribers[topic]...)
    var wg sync.WaitGroup
    for _, subscriber := range subscribers {
        wg.Add(1)
        go func(ch chan any) {
            defer wg.Done()
            ch <- event
        }(subscriber)
    }
    wg.Wait()
}


func (eb *EventBus) Subscribe(topic string) chan any {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan any)
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	return ch
}