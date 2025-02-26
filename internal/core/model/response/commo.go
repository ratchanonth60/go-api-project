package response

type ErrorResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e *ErrorResponse) Error() string {
	return e.Msg
}

type SuccResponse struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
