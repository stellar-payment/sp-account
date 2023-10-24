package indto

type UserParams struct {
	UserID     string
	ListUserID []string
	Username   string
	RoleID     int64
	Keyword    string
	Limit      uint64
	Page       uint64
}

type User struct {
	UserID   string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	RoleID   int64  `db:"role_id"`
}
