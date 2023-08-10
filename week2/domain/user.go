package domain

//User 领域对象，是DDD中的entity
//BO(business object)

type User struct {
	Id       int64
	Email    string
	Password string

	//添加如下字段，用户昵称，生日和个人简介
	Nickname string
	Birthday string
	Abstract string
	Ctime    string
	Utime    string
}
