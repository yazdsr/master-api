package request

import "time"

type CreateUser struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	FullName   string    `json:"full_name"`
	ServerID   int       `json:"server_id"`
	StartDate  time.Time `json:"start_date"`
	ValidUntil time.Time `json:"valid_until"`
}

type UpdateUser struct {
	Password   string    `json:"password"`
	FullName   string    `json:"full_name"`
	StartDate  time.Time `json:"start_date"`
	ValidUntil time.Time `json:"valid_until"`
}
