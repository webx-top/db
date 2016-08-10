package xorm

import (
	. "github.com/webx-top/dbx/driver"
)

type Select struct {
	table     string
	condition CondBuilder
	having    CondBuilder
	group     []string
	limit     int
	offset    int
	sort      []string
	distinct  []string
}

func (s *Select) SetTable(table string) Selecter {
	s.table = table
	return s
}

func (s *Select) SetCondition(builder CondBuilder) Selecter {
	s.condition = builder
	return s
}

func (s *Select) SetDistinct(fields ...string) Selecter {
	s.distinct = fields
	return s
}

func (s *Select) SetSort(fields ...string) Selecter {
	return s
}

func (s *Select) SetLimit(args ...int) Selecter {
	return s
}

func (s *Select) SetGroup(fields ...string) Selecter {
	return s

}

func (s *Select) SetHaving(builder CondBuilder) Selecter {
	return s
}

func (s *Select) Table() string {
	return s.table
}

func (s *Select) Condition() CondBuilder {
	return s.condition
}

func (s *Select) Distinct() []string {
	return s.distinct
}

func (s *Select) Sort() []string {
	return s.sort
}

func (s *Select) Limit() (int, int) {
	return s.limit, s.offset
}

func (s *Select) Group() []string {
	return s.group
}

func (s *Select) Having() CondBuilder {
	return s.having
}
