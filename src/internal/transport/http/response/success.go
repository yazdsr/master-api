package response

type Success struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
}