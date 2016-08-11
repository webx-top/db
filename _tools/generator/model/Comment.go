package model

type Comment struct {
	Id         	int64   	`db:"id"`
	Content    	string  	`db:"content"`
	Quote      	string  	`db:"quote"`
	Etype      	string  	`db:"etype"`
	RootId     	int64   	`db:"root_id"`
	RId        	int64   	`db:"r_id"`
	RType      	string  	`db:"r_type"`
	RelatedTimes	int     	`db:"related_times"`
	RootTimes  	int     	`db:"root_times"`
	Uid        	int64   	`db:"uid"`
	Uname      	string  	`db:"uname"`
	Up         	int64   	`db:"up"`
	Down       	int64   	`db:"down"`
	Created    	int     	`db:"created"`
	Updated    	int     	`db:"updated"`
	Status     	int     	`db:"status"`
	RcId       	int64   	`db:"rc_id"`
	RcType     	string  	`db:"rc_type"`
	ForUname   	string  	`db:"for_uname"`
	ForUid     	int64   	`db:"for_uid"`
}
