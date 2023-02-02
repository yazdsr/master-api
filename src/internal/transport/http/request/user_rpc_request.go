package request

type CreateUserRPC struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteUserRPC struct {
	Username string `json:"username"`
}