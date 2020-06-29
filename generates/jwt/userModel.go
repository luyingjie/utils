package jwt

type UserModel struct {
	UserId   int64
	UserKey  string
	UserName string
	PassWord string
	TimeUnix int64
}
