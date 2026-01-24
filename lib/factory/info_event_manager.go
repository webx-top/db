package factory

import "sync"

type EventManager interface {
	GetOk(table string) (*Event, bool)
	MustGetEvent(table string) *Event
	Range(f func(table string, evt *Event) bool)
}

func NewEventManager(unsafe bool) EventManager {
	if unsafe {
		return &UnsafeMapEvents{}
	}
	return &SafeMapEvents{m: map[string]*Event{}}
}

type SafeMapEvents struct {
	m  map[string]*Event
	mu sync.RWMutex
}

func (e *SafeMapEvents) GetOk(table string) (*Event, bool) {
	e.mu.RLock()
	evt, ok := e.m[table]
	e.mu.RUnlock()
	return evt, ok
}

func (e *SafeMapEvents) MustGetEvent(table string) *Event {
	evt, ok := e.GetOk(table)
	if ok {
		return evt
	}
	evt = NewEvent()
	e.mu.Lock()
	e.m[table] = evt
	e.mu.Unlock()
	return evt
}

func (e *SafeMapEvents) Range(f func(table string, evt *Event) bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for k, v := range e.m {
		if !f(k, v) {
			break
		}
	}
}

type UnsafeMapEvents map[string]*Event

func (e UnsafeMapEvents) GetOk(table string) (*Event, bool) {
	evt, ok := e[table]
	return evt, ok
}

func (e *UnsafeMapEvents) MustGetEvent(table string) *Event {
	evt, ok := e.GetOk(table)
	if ok {
		return evt
	}
	evt = NewEvent()
	(*e)[table] = evt
	return evt
}

func (e UnsafeMapEvents) Range(f func(table string, evt *Event) bool) {
	for k, v := range e {
		if !f(k, v) {
			break
		}
	}
}
