package rest_err

type RestErr interface {
	Error() string
	StatusCode() int
	Errors() []string
}
