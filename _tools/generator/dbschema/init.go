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

var AllfieldsMap FieldValidator=map[string]map[string]bool{"user":map[string]bool{"id":true, "passwd":true, "salt":true, "email":true, "created":true, "updated":true, "active":true, "avatar":true, "uname":true, "mobile":true, "login_time":true, "login_ip":true}, "attathment":map[string]bool{"name":true, "size":true, "uid":true, "deleted":true, "created":true, "rc_id":true, "rc_type":true, "id":true, "path":true, "extension":true, "type":true, "audited":true, "tags":true}, "category":map[string]bool{"name":true, "haschild":true, "updated":true, "rc_type":true, "sort":true, "tmpl":true, "id":true, "pid":true, "description":true}, "comment":map[string]bool{"root_id":true, "uid":true, "created":true, "rc_type":true, "r_type":true, "root_times":true, "rc_id":true, "quote":true, "r_id":true, "uname":true, "id":true, "content":true, "etype":true, "related_times":true, "up":true, "down":true, "updated":true, "status":true, "for_uname":true, "for_uid":true}, "link":map[string]bool{"url":true, "logo":true, "show":true, "created":true, "sort":true, "id":true, "name":true, "verified":true, "updated":true, "catid":true}, "ocontent":map[string]bool{"id":true, "rc_id":true, "rc_type":true, "content":true, "etype":true}, "post":map[string]bool{"description":true, "content":true, "created":true, "uid":true, "views":true, "allow_comment":true, "tags":true, "display":true, "comments":true, "deleted":true, "month":true, "catid":true, "id":true, "title":true, "etype":true, "updated":true, "uname":true, "passwd":true, "likes":true, "year":true}, "album":map[string]bool{"catid":true, "id":true, "deleted":true, "tags":true, "likes":true, "display":true, "title":true, "created":true, "updated":true, "comments":true, "allow_comment":true, "description":true, "content":true, "views":true}, "config":map[string]bool{"key":true, "val":true, "updated":true, "id":true}, "tag":map[string]bool{"uid":true, "created":true, "times":true, "rc_type":true, "id":true, "name":true}}
