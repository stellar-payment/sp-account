package model

type User struct {
	UserID   string  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	RoleID   int64  `db:"role_id"`
}
