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

var AllfieldsMap FieldValidator=map[string]map[string]bool{"config":map[string]bool{"id":true, "key":true, "val":true, "updated":true}, "link":map[string]bool{"id":true, "name":true, "verified":true, "catid":true, "sort":true, "url":true, "logo":true, "show":true, "created":true, "updated":true}, "post":map[string]bool{"content":true, "views":true, "catid":true, "deleted":true, "year":true, "id":true, "created":true, "updated":true, "uid":true, "passwd":true, "title":true, "likes":true, "month":true, "allow_comment":true, "tags":true, "description":true, "etype":true, "display":true, "uname":true, "comments":true}, "album":map[string]bool{"deleted":true, "display":true, "allow_comment":true, "likes":true, "title":true, "content":true, "created":true, "views":true, "tags":true, "catid":true, "id":true, "updated":true, "comments":true, "description":true}, "category":map[string]bool{"id":true, "pid":true, "name":true, "description":true, "haschild":true, "updated":true, "rc_type":true, "sort":true, "tmpl":true}, "ocontent":map[string]bool{"rc_type":true, "content":true, "etype":true, "id":true, "rc_id":true}, "tag":map[string]bool{"times":true, "rc_type":true, "id":true, "name":true, "uid":true, "created":true}, "user":map[string]bool{"login_time":true, "login_ip":true, "updated":true, "uname":true, "passwd":true, "salt":true, "mobile":true, "avatar":true, "id":true, "email":true, "created":true, "active":true}, "attathment":map[string]bool{"id":true, "name":true, "path":true, "uid":true, "deleted":true, "created":true, "rc_id":true, "extension":true, "type":true, "size":true, "audited":true, "rc_type":true, "tags":true}, "comment":map[string]bool{"etype":true, "related_times":true, "uid":true, "up":true, "id":true, "quote":true, "r_id":true, "r_type":true, "root_times":true, "root_id":true, "uname":true, "down":true, "created":true, "status":true, "rc_id":true, "rc_type":true, "for_uname":true, "content":true, "updated":true, "for_uid":true}}
