package response

type Err struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ErrorRPC struct {
	Jsonrpc string `json:"jsonrpc"`
	Error   Err    `json:"error"`
	ID      int    `json:"id"`
}

type Result string

type SuccessRPC struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
	ID      int    `json:"id"`
}
