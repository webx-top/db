// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package mongo wraps the gopkg.in/mgo.v2 MongoDB driver. See
// https://upper.io/db.v3/mongo for documentation, particularities and usage
// examples.
package mongo

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/webx-top/db"
	mgo "github.com/webx-top/qmgo"
	"github.com/webx-top/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter holds the name of the mongodb adapter.
const Adapter = `mongo`

var ConnTimeout = time.Second * 5

// Source represents a MongoDB database.
type Source struct {
	db.Settings

	name          string
	connURL       db.ConnectionURL
	client        *mgo.Client
	database      *mgo.Database
	version       []int
	collections   map[string]*Collection
	collectionsMu sync.Mutex
}

func init() {
	db.RegisterAdapter(Adapter, &db.AdapterFuncMap{
		Open: Open,
	})
}

// Open stablishes a new connection to a SQL server.
func Open(settings db.ConnectionURL) (db.Database, error) {
	d := &Source{Settings: db.NewSettings()}
	if err := d.Open(settings); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Source) ConnectionURL() db.ConnectionURL {
	return s.connURL
}

// SetConnMaxLifetime is not supported.
func (s *Source) SetConnMaxLifetime(time.Duration) {
	s.Settings.SetConnMaxLifetime(time.Duration(0))
}

// SetMaxIdleConns is not supported.
func (s *Source) SetMaxIdleConns(int) {
	s.Settings.SetMaxIdleConns(0)
}

// SetMaxOpenConns is not supported.
func (s *Source) SetMaxOpenConns(int) {
	s.Settings.SetMaxOpenConns(0)
}

// Name returns the name of the database.
func (s *Source) Name() string {
	return s.name
}

// Open attempts to connect to the database.
func (s *Source) Open(connURL db.ConnectionURL) error {
	s.connURL = connURL
	return s.open()
}

// Clone returns a cloned db.Database session.
func (s *Source) Clone() (db.Database, error) {
	newClient := *s.client
	clone := &Source{
		Settings: db.NewSettings(),

		name:        s.name,
		connURL:     s.connURL,
		client:      &newClient,
		database:    newClient.Database(s.database.GetDatabaseName()),
		version:     s.version,
		collections: map[string]*Collection{},
	}
	return clone, nil
}

// NewTransaction should support transactions, version of mongoDB server >= v4.0
func (s *Source) NewTransaction(ctx context.Context, opt ...*options.SessionOptions) (db.Tx, error) {
	if !s.versionAtLeast(4) {
		return nil, db.ErrUnsupported
	}
	sessionOpts := opts.Session()
	if len(opt) > 0 && opt[0].SessionOptions != nil {
		sessionOpts = opt[0].SessionOptions
	}
	session, err := s.client.Raw().StartSession(sessionOpts)
	if err != nil {
		return nil, err
	}
	return &txImpl{
		session: session,
		context: ctx,
	}, db.ErrUnsupported
}

func (s *Source) DoTransaction(ctx context.Context, callback func(sessCtx context.Context) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
	return s.client.DoTransaction(ctx, callback, opts...)
}

type txImpl struct {
	session mongo.Session
	context context.Context
}

// Rollback discards all the instructions on the current transaction.
func (t *txImpl) Rollback() error {
	defer t.session.EndSession(t.context)
	return t.session.AbortTransaction(t.context)
}

// Commit commits the current transaction.
func (t *txImpl) Commit() error {
	defer t.session.EndSession(t.context)
	return t.session.CommitTransaction(t.context)
}

func (t *txImpl) Context() context.Context {
	return t.context
}

// Ping checks whether a connection to the database is still alive by pinging
// it, establishing a connection if necessary.
func (s *Source) Ping() error {
	return s.client.Ping(ConnTimeout.Milliseconds())
}

func (s *Source) ClearCache() {
	s.collectionsMu.Lock()
	defer s.collectionsMu.Unlock()
	s.collections = make(map[string]*Collection)
}

// Driver returns the underlying *mgo.Client instance.
func (s *Source) Driver() interface{} {
	return s.client
}

func (s *Source) open() error {
	var err error

	ctx := context.Background()
	connTimeoutMs := ConnTimeout.Milliseconds()
	cfg := &mgo.Config{
		Uri:              s.connURL.String(),
		ConnectTimeoutMS: &connTimeoutMs,
		Auth: &mgo.Credential{
			AuthMechanism: `SCRAM-SHA-256`,
		},
		Database: s.connURL.(ConnectionURL).Database,
	}
	if s.client, err = mgo.NewClient(ctx, cfg); err != nil {
		return err
	}

	s.collections = map[string]*Collection{}
	if len(cfg.Database) > 0 {
		s.database = s.client.Database(cfg.Database)
	} else {
		s.database = s.client.Database(`test`)
	}

	return nil
}

// Close terminates the current database session.
func (s *Source) Close() error {
	if s.client != nil {
		s.client.Close(context.Background())
	}
	return nil
}

// Collections returns a list of non-system tables from the database.
func (s *Source) Collections() (cols []string, err error) {
	var rawcols []string
	var col string
	ctx := context.Background()
	if rawcols, err = s.database.Raw().ListCollectionNames(ctx, mgo.D{}); err != nil {
		return nil, err
	}

	cols = make([]string, 0, len(rawcols))

	for _, col = range rawcols {
		if !strings.HasPrefix(col, "system.") {
			cols = append(cols, col)
		}
	}

	return cols, nil
}

// Collection returns a collection by name.
func (s *Source) Collection(name string) db.Collection {
	s.collectionsMu.Lock()
	defer s.collectionsMu.Unlock()

	var col *Collection
	var ok bool

	if col, ok = s.collections[name]; !ok {
		col = &Collection{
			parent:     s,
			collection: s.database.Collection(name),
		}
		s.collections[name] = col
	}

	return col
}

func (s *Source) versionAtLeast(version ...int) bool {
	// only fetch this once - it makes a db call
	if len(s.version) == 0 {
		version := s.client.ServerVersion()
		vp := strings.SplitN(version, `.`, 4)
		versionArray := make([]int, len(vp))
		for i, v := range vp {
			n, e := strconv.Atoi(v)
			if e != nil {
				break
			}
			versionArray[i] = n
		}
		s.version = versionArray
	}

	// Check major version first
	if s.version[0] > version[0] {
		return true
	}

	for i := range version {
		if i == len(s.version) {
			return false
		}
		if s.version[i] < version[i] {
			return false
		}
	}
	return true
}
