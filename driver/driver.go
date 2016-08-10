package driver

import (
	"time"
)

type Driver interface {
	Connect(interface{}, string, ...time.Duration) error

	//Read
	All(Selecter, interface{}, ...string) error
	Count(Selecter, interface{}, ...string) (int, error)
	One(Selecter, interface{}, ...string) error

	//Write
	Delete(interface{}, CondBuilder, ...string) error
	Update(interface{}, H, CondBuilder, ...string) error
	Insert(interface{}, H, ...string) error
	Upsert(interface{}, H, CondBuilder, ...string) (int, error)
}

type H struct {
	data map[string]interface{}
	keys []string
}

func (h *H) Add(key string, value interface{}) {
	h.data[key] = value
	h.keys = append(h.keys, key)
}

func (h *H) Del(key string) {
	if _, ok := h.data[key]; ok {
		delete(h.data, key)
		keys := h.keys
		h.keys = []string{}
		for _, ekey := range keys {
			if ekey != key {
				h.keys = append(h.keys, ekey)
			}
		}
	}
}

func (h *H) Get(key string) interface{} {
	if v, ok := h.data[key]; ok {
		return v
	}
	return nil
}

func (h *H) Has(key string) bool {
	if _, ok := h.data[key]; ok {
		return true
	}
	return false
}

func (h *H) Map() map[string]interface{} {
	return h.data
}

func (h *H) Keys() []string {
	return h.keys
}

func (h *H) Build() interface{} {
	return h.data
}

type CondBuilder interface {
	Build() interface{}
}

type Selecter interface {
	//Setter
	SetTable(string) Selecter
	SetCondition(CondBuilder) Selecter
	SetDistinct(...string) Selecter
	SetSort(...string) Selecter
	SetLimit(...int) Selecter
	SetGroup(...string) Selecter
	SetHaving(CondBuilder) Selecter

	//Getter
	Table() string
	Condition() CondBuilder
	Distinct() []string
	Sort() []string
	Limit() (int, int) //limit,offset
	Group() []string
	Having() CondBuilder
}
