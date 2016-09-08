//Generated by webx-top/db
package dbschema

import (
	"github.com/webx-top/db/lib/factory"
)

func init(){
	factory.Fields=map[string]map[string]*factory.FieldInfo{"tag":map[string]*factory.FieldInfo{"name":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"标签名"}, "uid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"创建者"}, "created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"创建时间"}, "times":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"使用次数"}, "rc_type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"post", Comment:"关联类型"}, "id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"ID"}}, "user":map[string]*factory.FieldInfo{"uname":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"用户名"}, "passwd":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:64, Options:[]string{}, DefaultValue:"", Comment:"密码"}, "login_ip":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:40, Options:[]string{}, DefaultValue:"", Comment:"上次登录IP"}, "created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"创建时间"}, "avatar":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"头像"}, "active":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"Y", "N"}, DefaultValue:"Y", Comment:"激活"}, "id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"UID"}, "salt":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:64, Options:[]string{}, DefaultValue:"", Comment:"盐值"}, "email":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:100, Options:[]string{}, DefaultValue:"", Comment:"邮箱"}, "mobile":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:15, Options:[]string{}, DefaultValue:"", Comment:"手机号"}, "login_time":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"上次登录时间"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"更新时间"}}, "album":map[string]*factory.FieldInfo{"created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"创建时间"}, "views":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"浏览次数"}, "catid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"分类ID"}, "id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"ID"}, "display":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"ALL", "SELF", "FRIEND", "PWD"}, DefaultValue:"ALL", Comment:"显示"}, "allow_comment":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"Y", "N"}, DefaultValue:"Y", Comment:"是否允许评论"}, "title":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:180, Options:[]string{}, DefaultValue:"", Comment:"标题"}, "content":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{}, DefaultValue:"", Comment:"正文"}, "tags":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:255, Options:[]string{}, DefaultValue:"", Comment:"标签"}, "description":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"简介"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"编辑时间"}, "comments":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"评论次数"}, "likes":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"喜欢次数"}, "deleted":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"删除时间"}}, "attathment":map[string]*factory.FieldInfo{"name":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:100, Options:[]string{}, DefaultValue:"", Comment:"文件名"}, "path":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:255, Options:[]string{}, DefaultValue:"", Comment:"保存路径"}, "created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"创建时间"}, "audited":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"审核时间"}, "rc_id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"关联id"}, "deleted":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"被删除时间"}, "rc_type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"关联类型"}, "tags":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:255, Options:[]string{}, DefaultValue:"", Comment:"标签"}, "id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"ID"}, "extension":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:5, Options:[]string{}, DefaultValue:"", Comment:"扩展名"}, "type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"image", "media", "other"}, DefaultValue:"image", Comment:"文件类型"}, "size":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"", Comment:"文件尺寸"}, "uid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"UID"}}, "ocontent":map[string]*factory.FieldInfo{"etype":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"markdown"}, DefaultValue:"markdown", Comment:"编辑器类型"}, "id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"ID"}, "rc_id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"关联ID"}, "rc_type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"post", Comment:"关联类型"}, "content":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{}, DefaultValue:"", Comment:"博客原始内容"}}, "link":map[string]*factory.FieldInfo{"id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"主键ID"}, "sort":&factory.FieldInfo{DataType:"int", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"排序"}, "show":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"Y", "N"}, DefaultValue:"N", Comment:"是否显示"}, "verified":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"验证时间"}, "created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"创建时间"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"更新时间"}, "catid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"分类"}, "name":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"名称"}, "url":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"网址"}, "logo":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"LOGO"}}, "post":map[string]*factory.FieldInfo{"id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"ID"}, "created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"创建时间"}, "likes":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"被喜欢次数"}, "month":&factory.FieldInfo{DataType:"string", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:1, Options:[]string{}, DefaultValue:"", Comment:"归档月份"}, "title":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:180, Options:[]string{}, DefaultValue:"", Comment:"标题"}, "content":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{}, DefaultValue:"", Comment:"内容"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"修改时间"}, "display":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"ALL", "SELF", "FRIEND", "PWD"}, DefaultValue:"ALL", Comment:"显示"}, "uname":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"作者名"}, "views":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"被浏览次数"}, "comments":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"被评论次数"}, "deleted":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"被删除时间"}, "allow_comment":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"Y", "N"}, DefaultValue:"Y", Comment:"是否允许评论"}, "catid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"分类ID"}, "description":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"简介"}, "etype":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"html", "markdown"}, DefaultValue:"html", Comment:"编辑器类型"}, "uid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"UID"}, "passwd":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:64, Options:[]string{}, DefaultValue:"", Comment:"访问密码"}, "year":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:5, Options:[]string{}, DefaultValue:"", Comment:"归档年份"}, "tags":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:255, Options:[]string{}, DefaultValue:"", Comment:"标签"}}, "category":map[string]*factory.FieldInfo{"name":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"分类名称"}, "description":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"说明"}, "haschild":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"Y", "N"}, DefaultValue:"N", Comment:"是否有子分类"}, "rc_type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"post", Comment:"关联类型"}, "id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"ID"}, "pid":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"上级分类"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"更新时间"}, "sort":&factory.FieldInfo{DataType:"int", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"排序"}, "tmpl":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:100, Options:[]string{}, DefaultValue:"", Comment:"模板"}}, "comment":map[string]*factory.FieldInfo{"uid":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:"发布者id"}, "status":&factory.FieldInfo{DataType:"int", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:1, Options:[]string{}, DefaultValue:"0", Comment:"状态"}, "for_uname":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"被评人用户名"}, "for_uid":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:"被评人id"}, "content":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{}, DefaultValue:"", Comment:"内容"}, "r_type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{"reply", "append"}, DefaultValue:"", Comment:"关联类型（reply-回复，append-追加）"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"更新时间"}, "rc_id":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:"关联内容ID"}, "etype":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"html", Comment:"编辑器类型"}, "up":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:"被顶次数"}, "root_id":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:""}, "r_id":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:"关联本表的id"}, "down":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"0", Comment:"被踩次数"}, "id":&factory.FieldInfo{DataType:"int64", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:20, Options:[]string{}, DefaultValue:"", Comment:"主键"}, "quote":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:0, Options:[]string{}, DefaultValue:"", Comment:"引用内容"}, "uname":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"", Comment:"发布者用户名"}, "created":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"创建时间"}, "rc_type":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:30, Options:[]string{}, DefaultValue:"post", Comment:"关联内容类型"}, "related_times":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"本身被回复次数"}, "root_times":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"根节点下的所有回复次数"}}, "config":map[string]*factory.FieldInfo{"id":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:true, AutoIncrement:true, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"", Comment:"主键ID"}, "key":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:60, Options:[]string{}, DefaultValue:"", Comment:"配置项"}, "val":&factory.FieldInfo{DataType:"string", Unsigned:false, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:200, Options:[]string{}, DefaultValue:"", Comment:"配置值"}, "updated":&factory.FieldInfo{DataType:"int", Unsigned:true, PrimaryKey:false, AutoIncrement:false, Min:0, Max:0, MaxSize:10, Options:[]string{}, DefaultValue:"0", Comment:"更新时间"}}}

}

