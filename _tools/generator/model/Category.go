package model

type Category struct {
	Id      	int     	`db:"id"`
	Pid     	int     	`db:"pid"`
	Name    	string  	`db:"name"`
	Description	string  	`db:"description"`
	Haschild	string  	`db:"haschild"`
	Updated 	int     	`db:"updated"`
	RcType  	string  	`db:"rc_type"`
	Sort    	int     	`db:"sort"`
	Tmpl    	string  	`db:"tmpl"`
}
