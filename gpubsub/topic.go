package gpubsub

import "sync"

type Topic[T any] struct {
	mu            sync.RWMutex
	name          string
	subscriptions map[string]*Subscription[T]
}

func NewTopic[T any](name string) *Topic[T] {
	return &Topic[T]{
		name:          name,
		subscriptions: make(map[string]*Subscription[T]),
	}
}

func (t *Topic[T]) Name() string {
	return t.name
}

func (t *Topic[T]) Subscriptions() map[string]*Subscription[T] {
	return t.subscriptions
}

func (t *Topic[T]) Publish(message T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, subscription := range t.subscriptions {
		select {
		case subscription.ch <- message:
		default:
		}
	}
}

func (t *Topic[T]) NewSubscription(name string, bufferSize int, concurrency int64) *Subscription[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subscriptions[name]; !ok {
		t.subscriptions[name] = &Subscription[T]{
			name:        name,
			topic:       t,
			ch:          make(chan T, bufferSize),
			concurrency: concurrency,
		}
	}
	return t.subscriptions[name]
}

func (t *Topic[T]) Register(subscription *Subscription[T]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subscriptions[subscription.name]; !ok {
		t.subscriptions[subscription.name] = subscription
	}
}

func (t *Topic[T]) Unregister(subscription *Subscription[T]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subscriptions[subscription.name]; ok {
		delete(t.subscriptions, subscription.name)
	}
}
