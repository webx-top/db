package factory

import "github.com/webx-top/db"

func NewEvents() Events {
	return Events{}
}

type Events map[string]*Event

func (e Events) Call(event string, model Model, editColumns []string, mw func(db.Result) db.Result, args ...interface{}) error {
	if len(args) > 0 {
		table := model.Short_()
		events := []*Event{}
		if evt, ok := e[table]; ok {
			if evt.Exists(event) {
				events = append(events, evt)
			}
		}
		if evt, ok := e[`*`]; ok {
			if evt.Exists(event) {
				events = append(events, evt)
			}
		}
		if len(events) == 0 {
			return nil
		}
		rows := model.NewObjects()
		num := int64(1000)
		cnt, err := model.ListByOffset(&rows, mw, 0, int(num), args...)
		if err != nil {
			return err
		}
		total := cnt()
		if total < 1 {
			return nil
		}
		for i := int64(0); i < total; i += num {
			if i > 0 {
				rows = model.NewObjects()
				_, err := model.ListByOffset(&rows, mw, int(i), int(num), args...)
				if err != nil {
					return err
				}
			}
			return rows.Range(func(m Model) error {
				for _, evt := range events {
					if err := evt.Call(event, m, editColumns...); err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
	return e.call(event, model, editColumns...)
}

func (e Events) call(event string, model Model, editColumns ...string) error {
	table := model.Short_()
	if evt, ok := e[table]; ok {
		err := evt.Call(event, model, editColumns...)
		if err != nil {
			return err
		}
	}
	if evt, ok := e[`*`]; ok {
		return evt.Call(event, model, editColumns...)
	}
	return nil
}

func (e *Events) On(event string, h EventHandler, table string, async ...bool) {
	evt, ok := (*e)[table]
	if !ok {
		evt = NewEvent()
		(*e)[table] = evt
	}
	evt.On(event, h, async...)
}