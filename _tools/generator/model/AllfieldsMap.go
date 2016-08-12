package model
type FieldValidator map[string]map[string]bool

func (f FieldValidator) ValidField(table string,field string) bool {
	if tb,ok := f[table]; ok {
		return tb[field]
	}
	return false
}

func (f FieldValidator) ValidTable(table string) bool {
	_,ok := f[table]
	return ok
}

var AllfieldsMap FieldValidator=map[string]map[string]bool{"album":map[string]bool{"id":true, "title":true, "likes":true, "display":true, "allow_comment":true, "content":true, "tags":true, "description":true, "views":true, "comments":true, "deleted":true, "created":true, "updated":true, "catid":true}, "comment":map[string]bool{"root_id":true, "r_type":true, "up":true, "updated":true, "etype":true, "r_id":true, "related_times":true, "uid":true, "created":true, "for_uname":true, "id":true, "uname":true, "down":true, "rc_id":true, "rc_type":true, "content":true, "root_times":true, "status":true, "for_uid":true, "quote":true}, "config":map[string]bool{"key":true, "val":true, "updated":true, "id":true}, "ocontent":map[string]bool{"rc_type":true, "content":true, "etype":true, "id":true, "rc_id":true}, "post":map[string]bool{"uname":true, "views":true, "year":true, "month":true, "allow_comment":true, "catid":true, "etype":true, "uid":true, "display":true, "created":true, "updated":true, "tags":true, "description":true, "content":true, "passwd":true, "comments":true, "likes":true, "deleted":true, "id":true, "title":true}, "attathment":map[string]bool{"audited":true, "rc_type":true, "tags":true, "id":true, "uid":true, "created":true, "type":true, "size":true, "deleted":true, "rc_id":true, "name":true, "path":true, "extension":true}, "category":map[string]bool{"description":true, "haschild":true, "rc_type":true, "id":true, "pid":true, "name":true, "updated":true, "sort":true, "tmpl":true}, "link":map[string]bool{"catid":true, "id":true, "name":true, "logo":true, "created":true, "sort":true, "url":true, "show":true, "verified":true, "updated":true}, "tag":map[string]bool{"id":true, "name":true, "uid":true, "created":true, "times":true, "rc_type":true}, "user":map[string]bool{"uname":true, "email":true, "mobile":true, "login_time":true, "login_ip":true, "updated":true, "active":true, "id":true, "passwd":true, "salt":true, "created":true, "avatar":true}}
