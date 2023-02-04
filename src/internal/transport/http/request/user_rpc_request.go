package request

type CreateUserRPC struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteUserRPC struct {
	Username string `json:"username"`
}

type ActiveUserRPC struct {
	Username string `json:"username"`
}

type DisableUserRPC struct {
	Username string `json:"username"`
}

type UpdateUserRPC struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
