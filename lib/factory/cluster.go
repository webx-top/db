// added by swh@admpub.com
package factory

import (
	"log"
	"math/rand"

	"github.com/webx-top/db/lib/sqlbuilder"
)

// NewCluster : database cluster
func NewCluster() *Cluster {
	return &Cluster{
		masters: []sqlbuilder.Database{},
		slaves:  []sqlbuilder.Database{},
	}
}

// Cluster : database cluster
type Cluster struct {
	masters []sqlbuilder.Database
	slaves  []sqlbuilder.Database
	prefix  string
}

// W : write
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

// Prefix : table prefix
func (c *Cluster) Prefix() string {
	return c.prefix
}

// Table : Table full name (including the prefix)
func (c *Cluster) Table(tableName string) string {
	return c.prefix + tableName
}

// SetPrefix : setting table prefix
func (c *Cluster) SetPrefix(prefix string) {
	c.prefix = prefix
}

// R : read-only
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

// AddW : added writable database
func (c *Cluster) AddW(databases ...sqlbuilder.Database) {
	c.masters = append(c.masters, databases...)
}

// AddR : added read-only database
func (c *Cluster) AddR(databases ...sqlbuilder.Database) {
	c.slaves = append(c.slaves, databases...)
}

// SetW : set writable database
func (c *Cluster) SetW(index int, database sqlbuilder.Database) error {
	if len(c.masters) > index {
		c.masters[index] = database
		return nil
	}
	return ErrNotFoundKey
}

// SetR : set read-only database
func (c *Cluster) SetR(index int, database sqlbuilder.Database) error {
	if len(c.masters) > index {
		c.slaves[index] = database
		return nil
	}
	return ErrNotFoundKey
}

// CloseAll : Close all connections
func (c *Cluster) CloseAll() {
	c.CloseMasters()
	c.CloseSlaves()
}

// CloseMasters : Close all master connections
func (c *Cluster) CloseMasters() {
	for _, database := range c.masters {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

// CloseSlaves : Close all slave connections
func (c *Cluster) CloseSlaves() {
	for _, database := range c.slaves {
		if err := database.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

// CloseMaster : Close master connection
func (c *Cluster) CloseMaster(index int) bool {
	if len(c.masters) > index {
		if err := c.masters[index].Close(); err != nil {
			log.Println(err.Error())
		}
		return true
	}
	return false
}

// CloseSlave : Close slave connection
func (c *Cluster) CloseSlave(index int) bool {
	if len(c.slaves) > index {
		if err := c.slaves[index].Close(); err != nil {
			log.Println(err.Error())
		}
		return true
	}
	return false
}
