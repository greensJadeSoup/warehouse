package cp_error

import (
	"errors"
	"github.com/facebookarchive/stackerr"
	"warehouse/v5-go-component/cp_constant"
)

type SysError struct {
	Code int
	Message string
	se error
}

type NormalError struct {
	Code int
	Message string
	se error
}

func NewSysError(options ...interface{}) *SysError {
	cpErr := &SysError{
		Code: cp_constant.RESPONSE_CODE_SYSTEM,
	}

	for i, n := 0, len(options); i < n; i++ {
		switch err := options[i].(type) {
		case *SysError:
			cpErr = err
		case *NormalError:
			cpErr.Code = err.Code
			cpErr.Message = err.Message
			cpErr.se = err.se
		case string:
			cpErr.Message = err
			cpErr.se = stackerr.WrapSkip(errors.New(err), 1)
		case int:
			cpErr.Code = err
		case error:
			cpErr.Message = err.Error()
			cpErr.se = stackerr.WrapSkip(err, 1)
		default:
			cpErr.se = stackerr.New("invalid type of error")
		}
	}

	return cpErr
}

func NewNormalError(options ...interface{}) *NormalError {
	cpErr := &NormalError{
		Code: cp_constant.RESPONSE_CODE_COMMON_ERROR,
	}

	for i, n := 0, len(options); i < n; i++ {
		switch err := options[i].(type) {
		case *SysError:
			cpErr.Code = err.Code
			cpErr.Message = err.Message
			cpErr.se = err.se
		case *NormalError:
			cpErr = err
		case string:
			cpErr.Message = err
			cpErr.se = stackerr.WrapSkip(errors.New(err), 1)
		case int:
			cpErr.Code = err
		case error:
			cpErr.Message = err.Error()
			cpErr.se = stackerr.WrapSkip(err, 1)
		default:
			cpErr.se = stackerr.New("invalid type of error")
		}
	}

	return cpErr
}

func (this *SysError) Msg() string {
	return this.Message
}

func (this *NormalError) Msg() string {
	return this.Message
}

func (this *SysError) Error() string {
	return this.Message + "[Stack]: " + this.se.Error()
}

func (this *NormalError) Error() string {
	return this.Message + "[Stack]: " + this.se.Error()
}

func (this *SysError) Stack() []string {
	se, ok := interface{}(this.se).(*stackerr.Error)
	if ok {
		stacks := se.MultiStack().Stacks()
		ss := make([]string, 0)
		for _, v := range stacks {
			for _, vv := range v {
				ss = append(ss, vv.String())
			}
		}

		return ss
	}

	return []string{}
}

func (this *NormalError) Stack() []string {
	se, ok := interface{}(this.se).(*stackerr.Error)
	if ok {
		stacks := se.MultiStack().Stacks()
		ss := make([]string, 0)
		for _, v := range stacks {
			for _, vv := range v {
				ss = append(ss, vv.String())
			}
		}

		return ss
	}

	return []string{}
}

func (this *SysError) StackString() string {
	return "[Stack]: " + this.se.Error()
}

func (this *NormalError) StackString() string {
	return "[Stack]: " + this.se.Error()
}
