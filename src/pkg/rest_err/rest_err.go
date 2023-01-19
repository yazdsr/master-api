package rest_err

type restErr struct {
	statusCode int
	err        string
	errors     []string
}

func NewRestErr(statusCode int, err string, errors []string) RestErr {
	return &restErr{statusCode: statusCode, err: err, errors: errors}
}

func (re *restErr) Error() string {
	return re.err
}

func (re *restErr) StatusCode() int {
	return re.statusCode
}

func (re *restErr) Errors() []string {
	return re.errors
}
