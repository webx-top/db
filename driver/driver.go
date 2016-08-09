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
	Delete(string, CondBuilder, ...string) error
	Update(string, H, CondBuilder, ...string) error
	Insert(string, H, ...string) error
	Upsert(string, H, CondBuilder, ...string) (int, error)
}

type H map[string]interface{}

func (h H) Build() interface{} {
	return h
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
