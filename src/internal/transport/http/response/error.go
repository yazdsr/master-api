package response

type (
	Error struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Errors  []string `json:"errors,omitempty"`
	}
)
