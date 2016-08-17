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

var AllfieldsMap FieldValidator=map[string]map[string]bool{"album":map[string]bool{"description":true, "content":true, "updated":true, "display":true, "tags":true, "title":true, "views":true, "comments":true, "likes":true, "deleted":true, "catid":true, "id":true, "created":true, "allow_comment":true}, "attathment":map[string]bool{"uid":true, "deleted":true, "created":true, "name":true, "path":true, "extension":true, "type":true, "size":true, "audited":true, "rc_id":true, "rc_type":true, "id":true, "tags":true}, "ocontent":map[string]bool{"etype":true, "id":true, "rc_id":true, "rc_type":true, "content":true}, "link":map[string]bool{"verified":true, "sort":true, "url":true, "logo":true, "show":true, "updated":true, "catid":true, "id":true, "name":true, "created":true}, "post":map[string]bool{"views":true, "allow_comment":true, "id":true, "title":true, "updated":true, "created":true, "display":true, "uid":true, "likes":true, "deleted":true, "description":true, "content":true, "etype":true, "comments":true, "year":true, "catid":true, "tags":true, "uname":true, "passwd":true, "month":true}, "tag":map[string]bool{"created":true, "times":true, "rc_type":true, "id":true, "name":true, "uid":true}, "user":map[string]bool{"active":true, "avatar":true, "passwd":true, "salt":true, "login_time":true, "created":true, "login_ip":true, "updated":true, "id":true, "uname":true, "email":true, "mobile":true}, "category":map[string]bool{"haschild":true, "rc_type":true, "sort":true, "tmpl":true, "description":true, "pid":true, "name":true, "updated":true, "id":true}, "comment":map[string]bool{"etype":true, "r_type":true, "related_times":true, "up":true, "down":true, "updated":true, "rc_type":true, "content":true, "for_uid":true, "id":true, "r_id":true, "root_times":true, "uname":true, "rc_id":true, "for_uname":true, "quote":true, "uid":true, "created":true, "status":true, "root_id":true}, "config":map[string]bool{"id":true, "key":true, "val":true, "updated":true}}
