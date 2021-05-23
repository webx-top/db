// Copyright (c) 2021-present The webx-top/db authors. All rights reserved.
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

package clickhouse

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// ConnectionURL implements a MySQL connection struct.
type ConnectionURL struct {
	User     string
	Password string
	Database string
	Host     string

	// Optional
	TLSConfig              string
	NoDelay                bool
	Secure                 bool
	SkipVerify             bool
	Timeout                float64
	ReadTimeout            float64
	WriteTimeout           float64
	BlockSize              int64
	PoolSize               int64
	AltHosts               string
	ConnectionOpenStrategy string //random/in_order/time_random,
	Compress               bool
	Debug                  bool
	Options                map[string]string
}

func (c ConnectionURL) String() (s string) {
	// Adding protocol and address
	if len(c.Host) > 0 {
		host, port, err := net.SplitHostPort(c.Host)
		if err != nil {
			host = c.Host
			port = "9000"
		}
		s = s + fmt.Sprintf("tcp://%s:%s", host, port)
	} else {
		s = "tcp://127.0.0.1:9000"
	}

	// Do we have any options?
	if c.Options == nil {
		c.Options = map[string]string{}
	}

	// Converting options into URL values.
	vv := url.Values{
		`database`: []string{c.Database},
	}
	// Adding username.
	if len(c.User) > 0 {
		vv.Set(`username`, c.User)
		// Adding password.
		if len(c.Password) > 0 {
			vv.Set(`password`, c.Password)
		}
	}
	if len(c.TLSConfig) > 0 {
		vv.Set(`tls_config`, c.TLSConfig)
	}
	if c.NoDelay {
		vv.Set(`no_delay`, `true`)
	}
	if c.Secure {
		vv.Set(`secure`, `true`)
	}
	if c.SkipVerify {
		vv.Set(`skip_verify`, `true`)
	}
	if c.Timeout > 0 {
		vv.Set(`timeout`, fmt.Sprintf("%f", c.Timeout))
	}
	if c.ReadTimeout > 0 {
		vv.Set(`read_timeout`, fmt.Sprintf("%f", c.ReadTimeout))
	}
	if c.WriteTimeout > 0 {
		vv.Set(`write_timeout`, fmt.Sprintf("%f", c.WriteTimeout))
	}
	if c.BlockSize > 0 {
		vv.Set(`block_size`, fmt.Sprintf("%d", c.BlockSize))
	}
	if c.PoolSize > 0 {
		vv.Set(`pool_size`, fmt.Sprintf("%d", c.PoolSize))
	}
	if len(c.AltHosts) > 0 {
		vv.Set(`alt_hosts`, c.AltHosts)
	}
	if len(c.ConnectionOpenStrategy) > 0 {
		vv.Set(`connection_open_strategy`, c.ConnectionOpenStrategy)
	}
	if c.Compress {
		vv.Set(`compress`, `true`)
	}
	if c.Debug {
		vv.Set(`debug`, `true`)
	}
	// tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000

	for k, v := range c.Options {
		vv.Set(k, v)
	}

	// Inserting options.
	if p := vv.Encode(); len(p) > 0 {
		s = s + "?" + p
	}

	return s
}

// ParseURL parses s into a ConnectionURL struct.
func ParseURL(s string) (conn ConnectionURL, err error) {
	var url *url.URL
	if url, err = url.Parse(s); err != nil {
		return
	}
	query := url.Query()
	conn.User = query.Get(`username`)
	conn.Password = query.Get(`password`)
	conn.Host = url.Host
	conn.Database = query.Get(`database`)
	if v := query.Get(`tls_config`); len(v) > 0 {
		conn.TLSConfig = v
	}
	query.Del(`tls_config`)
	if v, e := strconv.ParseBool(query.Get(`no_delay`)); e == nil {
		conn.NoDelay = v
	}
	query.Del(`no_delay`)
	if v, e := strconv.ParseBool(query.Get(`secure`)); e == nil {
		conn.Secure = v
	}
	query.Del(`secure`)
	if v, e := strconv.ParseBool(query.Get(`skip_verify`)); e == nil {
		conn.SkipVerify = v
	}
	query.Del(`skip_verify`)
	if v, e := strconv.ParseFloat(query.Get(`timeout`), 64); e == nil {
		conn.Timeout = v
	}
	query.Del(`timeout`)
	if v, e := strconv.ParseFloat(query.Get(`read_timeout`), 64); e == nil {
		conn.ReadTimeout = v
	}
	query.Del(`read_timeout`)
	if v, e := strconv.ParseFloat(query.Get(`write_timeout`), 64); e == nil {
		conn.WriteTimeout = v
	}
	query.Del(`write_timeout`)
	if v, e := strconv.ParseInt(query.Get(`block_size`), 10, 64); e == nil {
		conn.BlockSize = v
	}
	query.Del(`block_size`)
	if v, e := strconv.ParseInt(query.Get(`pool_size`), 10, 64); e == nil {
		conn.PoolSize = v
	}
	query.Del(`pool_size`)
	if v := query.Get(`alt_hosts`); len(v) > 0 {
		conn.AltHosts = v
	}
	query.Del(`alt_hosts`)
	if v := query.Get(`connection_open_strategy`); len(v) > 0 {
		conn.ConnectionOpenStrategy = v
	}
	query.Del(`connection_open_strategy`)
	if v, e := strconv.ParseBool(query.Get(`compress`)); e == nil {
		conn.Compress = v
	}
	query.Del(`compress`)
	if v, e := strconv.ParseBool(query.Get(`debug`)); e == nil {
		conn.Debug = v
	}
	query.Del(`debug`)
	conn.Options = map[string]string{}

	for k := range query {
		conn.Options[k] = query.Get(k)
	}

	return
}
