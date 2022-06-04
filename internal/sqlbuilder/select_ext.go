package sqlbuilder

import (
	"github.com/webx-top/db"
	"github.com/webx-top/db/internal/sqladapter/exql"
)

// ForceIndex [SWH|+]
func (sel *selector) ForceIndex(index string) db.Selector {
	if len(index) == 0 {
		return sel
	}
	return sel.frame(func(sq *selectorQuery) error {
		sq.forceIndex = exql.JoinColumns(exql.Fragment(exql.ColumnWithName(index)))
		return nil
	})
}

func (sel *selector) RelationMap(relationMap map[string]BuilderChainFunc) db.Selector {
	return sel.frame(func(sq *selectorQuery) error {
		sq.relationMap = relationMap
		return nil
	})
}
