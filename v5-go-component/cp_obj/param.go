package cp_obj

import (
	"strings"
	"strconv"
	"errors"
)


//Param 请求参数集合
type Param map[string]string

//NewParam 创建新的请求参数集合
func NewParam() *Param {
	var pm Param = make(map[string]string)
	return &pm
}

//Map 获取Param为map类型
func (p Param) Map() map[string]string {
	return map[string]string(p)
}

//Add 添加请求参数
func (p Param) Add(key, val string) {
	key = strings.ToLower(key)
	p[key] = val
}

//GetString 获取请求参数：转为string，不存在时引发 error
func (p Param) GetString(key string) (string) {
	val, ok := p[key]
	if !ok {
		return ""
	}

	return val
}

//MustGetString 必须获取请求参数：转为string，不存在时返回零值或def
func (p Param) MustGetString(key string, def ...string) string {
	val := p.GetString(key)
	if val == "" && len(def) > 0 {
		return def[0]
	}

	return val
}

//GetInt 获取请求参数：转为int，不存在或格式错误时引发 error
func (p Param) GetInt(key string) (int, error) {
	val := p.GetString(key)
	if val == "" {
		return 0, errors.New(key + "参数格式为空")
	}

	i, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		return 0, errors.New(key + "参数格式错误")
	}

	return int(i), nil
}

//MustGetInt 必须获取请求参数：转为int，不存在或格式转化错误时返回零值或def
func (p Param) MustGetInt(key string, def ...int) int {
	i, err := p.GetInt(key)
	if err != nil && len(def) > 0 {
		return def[0]
	}

	return i
}

//GetInt64 获取请求参数：转为int64，不存在或格式错误时引发 error
func (p Param) GetInt64(key string) (int64, error) {
	val := p.GetString(key)
	if val == "" {
		return 0, errors.New(key + "参数格式为空")
	}

	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.New(key + "参数格式错误")
	}

	return i64, nil
}

//MustGetInt64 必须获取请求参数：转为int64，不存在或格式转化错误时返回零值或def
func (p Param) MustGetInt64(key string, def ...int64) int64 {
	i64, err := p.GetInt64(key)
	if err != nil && len(def) > 0 {
		return def[0]
	}

	return i64
}

//GetFloat32 获取请求参数：转为float32，不存在或格式错误时引发 error
func (p Param) GetFloat32(key string) (float32, error) {
	val := p.GetString(key)
	if val == "" {
		return 0, errors.New(key + "参数格式为空")
	}

	f32, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, errors.New(key + "参数格式错误")
	}

	return float32(f32), nil
}

//MustGetFloat32 必须获取请求参数：转为float32，不存在或格式转化错误时返回零值或def
func (p Param) MustGetFloat32(key string, def ...float32) float32 {
	f32, err := p.GetFloat32(key)
	if err != nil && len(def) > 0 {
		return def[0]
	}

	return f32
}

//GetFloat64 获取请求参数：转为float64，不存在或格式错误时引发 error
func (p Param) GetFloat64(key string) (float64, error) {
	val := p.GetString(key)
	if val == "" {
		return 0, errors.New(key + "参数格式为空")
	}

	f64, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.New(key + "参数格式错误")
	}

	return f64, nil
}

//MustGetFloat64 必须获取请求参数：转为float64，不存在或格式转化错误时返回零值或def
func (p Param) MustGetFloat64(key string, def ...float64) float64 {
	f64, err := p.GetFloat64(key)
	if err != nil && len(def) > 0 {
		return def[0]
	}

	return f64
}

//GetBool 获取请求参数：转为bool，不存在或格式错误时引发 error
func (p Param) GetBool(key string) (bool, error) {
	val := p.GetString(key)
	if val == "" {
		return false, errors.New(key + "参数格式为空")
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, errors.New(key + "参数格式错误")
	}

	return b, nil
}

//MustGetBool 必须获取请求参数：转为bool，不存在或格式转化错误时返回零值或def
func (p Param) MustGetBool(key string, def ...bool) bool {
	b, err := p.GetBool(key)
	if err != nil && len(def) > 0 {
		return def[0]
	}

	return b
}

//GetDatetime 获取请求参数：转为Datetime，不存在或格式错误时引发 error
func (p Param) GetDatetime(key string) (Datetime, error) {
	var dt Datetime

	val := p.GetString(key)
	if val == "" {
		return dt, errors.New(key + "参数格式为空")
	}

	err := dt.SetString(val)
	if err != nil {
		return dt, err
	}

	return dt, nil
}

//MustGetDatetime 必须获取请求参数：转为Datetime，不存在或格式转化错误时返回零值或def
func (p Param) MustGetDatetime(key string, def ...Datetime) Datetime {
	dt, err := p.GetDatetime(key)
	if err != nil && len(def) > 0 {
		return def[0]
	}

	return dt
}

//Del 删除请求参数
func (p Param) Del(key string) {
	key = strings.ToLower(key)
	delete(p, key)
}

//Clear 清空请求参数
func (p Param) Clear() {
	for k := range p {
		delete(p, k)
	}
}

//Has 是否存在请求参数
func (p Param) Has(key string) bool {
	_, ok := p[key]
	return ok
}
