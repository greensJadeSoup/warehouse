package cp_obj

import (
	"warehouse/v5-go-component/cp_constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
)

type Response struct {
	Code		int			`json:"code"`
	Message		string			`json:"message"`
	Stack		[]string		`json:"stack,omitempty"`
	Data		interface{}		`json:"data"`
}

func NewResponse() *Response {
	return &Response{}
}

func SpileResponse(err error) *Response {
	errNormal, ok := err.(*cp_error.NormalError)
	if ok {
		return &Response{Code: errNormal.Code, Message: errNormal.Msg(), Stack: errNormal.Stack()}
	}

	errSysError, ok := err.(*cp_error.SysError)
	if ok {
		return &Response{Code: errSysError.Code, Message: errSysError.Msg(), Stack: errSysError.Stack()}
	}

	return nil
}

func (this *Response) Err(options ...interface{}) *Response {
	this.Code = cp_constant.RESPONSE_CODE_COMMON_ERROR

	for i, n := 0, len(options); i < n; i++ {
		errNormal, ok := options[i].(*cp_error.NormalError)
		if ok {
			this.Code = errNormal.Code
			this.Message = errNormal.Msg()
			this.Stack = errNormal.Stack()
			cp_log.Warning(errNormal.StackString())
			continue
		}

		errSysError, ok := options[i].(*cp_error.SysError)
		if ok {
			this.Code = errSysError.Code
			this.Message = errSysError.Msg()
			this.Stack = errSysError.Stack()
			cp_log.Warning(errSysError.StackString())
			continue
		}

		switch opt := options[i].(type) {
		case error:
			this.Message = opt.Error()
		case int:
			this.Code = opt
		case string:
			this.Message = opt
		}
	}

	return this
}

func (this *Response) Ok(options... interface{}) *Response {
	this.Code = cp_constant.RESPONSE_CODE_OK
	if len(options) > 0 {
		this.Data = options[0]
	}

	if len(options) == 1 {
		return this
	}

	for i, n := 0, len(options); i < n; i++ {
		switch opt := options[i].(type) {
		case int:
			this.Code = opt
		case string:
			this.Message = opt
		}
	}

	return this
}
