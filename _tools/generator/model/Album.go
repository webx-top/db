package model

type Album struct {
	Id              	int     	`db:"id"`
	Title           	string  	`db:"title"`
	Description     	string  	`db:"description"`
	Content         	string  	`db:"content"`
	Created         	int     	`db:"created"`
	Updated         	int     	`db:"updated"`
	Views           	int     	`db:"views"`
	Comments        	int     	`db:"comments"`
	Likes           	int     	`db:"likes"`
	Display         	string  	`db:"display"`
	Deleted         	int     	`db:"deleted"`
	AllowComment    	string  	`db:"allow_comment"`
	Tags            	string  	`db:"tags"`
	Catid           	int     	`db:"catid"`
}
