// added by swh@admpub.com
package factory

import (
	"log"
	"math/rand"

	"github.com/webx-top/db/lib/sqlbuilder"
)

func NewCluster() *Cluster {
	return &Cluster{
		masters: []sqlbuilder.Database{},
		slaves:  []sqlbuilder.Database{},
	}
}

type Cluster struct {
	masters []sqlbuilder.Database
	slaves  []sqlbuilder.Database
	prefix  string
}

func (c *Cluster) W() sqlbuilder.Database {
	length := len(c.masters)
	if length == 0 {
		panic(`Not connected to any database`)
	}
	if length > 1 {
		return c.masters[rand.Intn(length-1)]
	}
	return c.masters[0]
}

func (c *Cluster) Prefix() string {
	return c.prefix
}

func (c *Cluster) Table(tableName string) string {
	return c.prefix + tableName
}

func (c *Cluster) SetPrefix(prefix string) {
	c.prefix = prefix
}

func (c *Cluster) R() sqlbuilder.Database {
	length := len(c.slaves)
	if length == 0 {
		return c.W()
	}
	if length > 1 {
		return c.slaves[rand.Intn(length-1)]
	}
	return c.slaves[0]
}

func (c *Cluster) AddW(databases ...sqlbuilder.Database) {
	c.masters = append(c.masters, databases...)
}

func (c *Cluster) AddR(databases ...sqlbuilder.Database) {
	c.slaves = append(c.slaves, databases...)
}

func (c *Cluster) SetW(index int, database sqlbuilder.Database) error {
	if len(c.masters) > index {
		c.masters[index] = database
		return nil
	}
	return ErrNotFoundKey
}

func (c *Cluster) SetR(index int, database sqlbuilder.Database) error {
	if len(c.masters) > index {
		c.slaves[index] = database
		return nil
	}
	return ErrNotFoundKey
}

func (c *Cluster) CloseAll() {
	c.CloseMasters()
	c.CloseSlaves()
}

func (c *Cluster) CloseMasters() {
	for _, database := range c.masters {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

func (c *Cluster) CloseSlaves() {
	for _, database := range c.slaves {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

func (c *Cluster) CloseMaster(index int) bool {
	if len(c.masters) > index {
		if err := c.masters[index].Close(); err != nil {
			log.Println(err.Error())
		}
		return true
	}
	return false
}

func (c *Cluster) CloseSlave(index int) bool {
	if len(c.slaves) > index {
		if err := c.slaves[index].Close(); err != nil {
			log.Println(err.Error())
		}
		return true
	}
	return false
}
