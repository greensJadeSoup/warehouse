package cp_obj

import (
	"warehouse/v5-go-component/cp_constant"
)

//easyjson:json
type TraceAction struct {
	LogLevel	cp_constant.TracingLevel `json:"log_level"`

	IP		string			`json:"ip,omitempty"`
	AppID		string			`json:"app_id,omitempty"`
	SessionKey	string			`json:"sessionkey,omitempty"`

	RequestID	string			`json:"request_id"`
	ChainID		string			`json:"chain_id"`
	ChainLevel	string			`json:"chain_level,omitempty"`
	ErrMsg		string			`json:"msg,omitempty"`

	//TimeStamp	int64			`json:"timestamp"`
	StartTime	string			`json:"start_time,omitempty"`
	TotalTime	int64			`json:"total_use_us"`

	Uri		string			`json:"uri,omitempty"`

	Body		string			`json:"body,omitempty"`
	Response 	*Response		`json:"response_data,omitempty"`
}

//easyjson:json
type TraceRuntime struct {
	LogLevel	cp_constant.TracingLevel `json:"log_level"`

	SvrName		string			`json:"svr_name"`
	RequestID	string			`json:"request_id"`
	ErrMsg		string			`json:"msg,omitempty"`
	StartTime	string			`json:"start_time,omitempty"`
}




