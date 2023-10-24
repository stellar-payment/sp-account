package dto

type UsersQueryParams struct {
	UserID  string `param:"userID"`
	RoleID  int64  `query:"role"`
	Keyword string `query:"keyword"`
	Limit   uint64 `query:"limit"`
	Page    uint64 `query:"page"`
}

type UserPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleID   int64  `json:"role_id" validate:"required"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
}

type ListUserResponse struct {
	Users []*UserResponse `json:"users"`
	Meta  ListPaginations `json:"meta"`
}
