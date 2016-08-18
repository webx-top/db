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

var AllfieldsMap FieldValidator=map[string]map[string]bool{"category":map[string]bool{"updated":true, "rc_type":true, "sort":true, "id":true, "name":true, "description":true, "haschild":true, "tmpl":true, "pid":true}, "ocontent":map[string]bool{"content":true, "etype":true, "id":true, "rc_id":true, "rc_type":true}, "post":map[string]bool{"id":true, "created":true, "year":true, "display":true, "uid":true, "passwd":true, "views":true, "deleted":true, "catid":true, "title":true, "etype":true, "comments":true, "likes":true, "month":true, "description":true, "content":true, "updated":true, "uname":true, "allow_comment":true, "tags":true}, "tag":map[string]bool{"uid":true, "created":true, "times":true, "rc_type":true, "id":true, "name":true}, "user":map[string]bool{"mobile":true, "login_ip":true, "updated":true, "uname":true, "email":true, "salt":true, "login_time":true, "created":true, "active":true, "avatar":true, "id":true, "passwd":true}, "album":map[string]bool{"description":true, "updated":true, "comments":true, "id":true, "created":true, "deleted":true, "allow_comment":true, "content":true, "catid":true, "likes":true, "views":true, "display":true, "tags":true, "title":true}, "attathment":map[string]bool{"name":true, "path":true, "size":true, "uid":true, "created":true, "audited":true, "rc_id":true, "id":true, "extension":true, "type":true, "deleted":true, "rc_type":true, "tags":true}, "comment":map[string]bool{"related_times":true, "rc_id":true, "for_uid":true, "id":true, "content":true, "r_type":true, "created":true, "status":true, "rc_type":true, "for_uname":true, "quote":true, "etype":true, "uid":true, "up":true, "uname":true, "down":true, "updated":true, "root_id":true, "r_id":true, "root_times":true}, "config":map[string]bool{"id":true, "key":true, "val":true, "updated":true}, "link":map[string]bool{"name":true, "url":true, "show":true, "sort":true, "catid":true, "id":true, "logo":true, "verified":true, "created":true, "updated":true}}
