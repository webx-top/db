package config

import (
	"net/url"
	"strings"
)

type Config struct {
	Engine  string
	User    string
	Pass    string
	Name    string
	Host    string
	Port    string
	Charset string
	Prefix  string
}

func (a *Config) String() string {
	var dsn string
	switch a.Engine {
	case `mysql`:
		var host string
		if strings.HasPrefix(a.Host, `unix:`) {
			host = "unix(" + strings.TrimPrefix(a.Host, `unix:`) + ")"
		} else {
			if a.Port == `` {
				a.Port = "3306"
			}
			host = "tcp(" + a.Host + ":" + a.Port + ")"
		}

		dsn = url.QueryEscape(a.User) + ":" + url.QueryEscape(a.Pass) + "@" + host + "/" + a.Name + "?charset=" + a.Charset
	case `mymysql`: //tcp:localhost:3306*gotest/root/root
		var host string
		if strings.HasPrefix(a.Host, `unix:`) {
			host = a.Host
		} else {
			if a.Port == `` {
				a.Port = "3306"
			}
			host = "tcp:" + a.Host + ":" + a.Port
		}
		dsn = host + "*" + a.Name + "/" + url.QueryEscape(a.User) + "/" + url.QueryEscape(a.Pass)
	default:
		panic(a.Engine + ` is not supported.`)
	}
	return dsn
}
