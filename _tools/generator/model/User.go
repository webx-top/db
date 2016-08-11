package model

type User struct {
	Id      	int     	`db:"id"`
	Uname   	string  	`db:"uname"`
	Passwd  	string  	`db:"passwd"`
	Salt    	string  	`db:"salt"`
	Email   	string  	`db:"email"`
	Mobile  	string  	`db:"mobile"`
	LoginTime	int     	`db:"login_time"`
	LoginIp 	string  	`db:"login_ip"`
	Created 	int     	`db:"created"`
	Updated 	int     	`db:"updated"`
	Active  	string  	`db:"active"`
	Avatar  	string  	`db:"avatar"`
}
