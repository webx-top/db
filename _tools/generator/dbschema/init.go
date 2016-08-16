package dbschema

import (
	"github.com/webx-top/db/lib/factory"
)

var Factory *factory.Factory = factory.DefaultFactory

type FieldValidator map[string]map[string]bool

func (f FieldValidator) ValidField(table string, field string) bool {
	if tb, ok := f[table]; ok {
		return tb[field]
	}
	return false
}

func (f FieldValidator) ValidTable(table string) bool {
	_, ok := f[table]
	return ok
}

var AllfieldsMap FieldValidator=map[string]map[string]bool{"album":map[string]bool{"description":true, "comments":true, "views":true, "catid":true, "id":true, "title":true, "content":true, "display":true, "created":true, "updated":true, "likes":true, "deleted":true, "allow_comment":true, "tags":true}, "category":map[string]bool{"haschild":true, "tmpl":true, "id":true, "name":true, "description":true, "sort":true, "pid":true, "updated":true, "rc_type":true}, "config":map[string]bool{"id":true, "key":true, "val":true, "updated":true}, "ocontent":map[string]bool{"etype":true, "id":true, "rc_id":true, "rc_type":true, "content":true}, "post":map[string]bool{"id":true, "uid":true, "uname":true, "year":true, "tags":true, "description":true, "updated":true, "views":true, "likes":true, "month":true, "comments":true, "deleted":true, "catid":true, "allow_comment":true, "title":true, "content":true, "etype":true, "created":true, "display":true, "passwd":true}, "attathment":map[string]bool{"name":true, "audited":true, "rc_id":true, "created":true, "id":true, "path":true, "extension":true, "type":true, "size":true, "uid":true, "deleted":true, "rc_type":true, "tags":true}, "comment":map[string]bool{"id":true, "related_times":true, "uname":true, "up":true, "updated":true, "r_id":true, "root_times":true, "down":true, "rc_type":true, "for_uname":true, "etype":true, "root_id":true, "r_type":true, "uid":true, "content":true, "quote":true, "created":true, "status":true, "rc_id":true, "for_uid":true}, "link":map[string]bool{"id":true, "url":true, "logo":true, "show":true, "catid":true, "name":true, "verified":true, "created":true, "updated":true, "sort":true}, "tag":map[string]bool{"id":true, "name":true, "uid":true, "created":true, "times":true, "rc_type":true}, "user":map[string]bool{"uname":true, "salt":true, "email":true, "updated":true, "avatar":true, "id":true, "passwd":true, "mobile":true, "login_time":true, "login_ip":true, "created":true, "active":true}}
