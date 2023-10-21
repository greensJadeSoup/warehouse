package cp_obj

import (
	"strings"
	"time"
)

type Datetime time.Time

//NewDatetime 创建一个当前时间的Datetime
func NewDatetime() Datetime {
	return Datetime(time.Now())
}

func (t *Datetime) UnmarshalJSON(data []byte) error {
	str := strings.Replace(string(data), `"`, "", -1)
	return t.SetString(str)
}

func (t Datetime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

//String 获取日期时间，格式：yyyy-MM-dd HH:mm:ss
func (t Datetime) String() string {
	return formatTimeToString(t.Time())
}

//SetString 设置日期时间，字符串格式：yyyy-MM-dd HH:mm:ss
func (t *Datetime) SetString(str string) error {
	tt, err := parseStringToTime(str)
	if err != nil {
		return err
	}

	*t = Datetime(tt)
	return nil
}

//Time 获取 Datetime的 time.Time 类型形式
func (t Datetime) Time() time.Time {
	return time.Time(t)
}

//SetTime 用 time.Time 类型来设置 Datetime 的值
func (t *Datetime) SetTime(tt time.Time) {
	*t = Datetime(tt)
}

//Unix 获取 Datetime 的 Unix 时间戮，单位：秒
func (t *Datetime) Unix() int64 {
	return t.Time().Unix()
}

//SetUnix 用时间戮来设置 Datetime 的值
func (t *Datetime) SetUnix(i64 int64) {
	*t = Datetime(time.Unix(i64, 0))
}
var timeFormat = "2006-01-02 15:04:05"

//将 yyyy-MM-dd HH:mm:ss 转换为 time.Time 类型
func parseStringToTime(str string) (time.Time, error) {
	var (
		loc *time.Location
		t   time.Time
		err error
	)
	if loc, err = time.LoadLocation("Local"); err == nil {
		if t, err = time.ParseInLocation(timeFormat, str, loc); err == nil {
			return t, nil
		}
	}
	return time.Now(), err
}

//将 time.Time 类型转换为 yyyy-MM-dd HH:mm:ss
func formatTimeToString(t time.Time) string {
	return t.Format(timeFormat)
}

func (t *Datetime) AddDate(y int, m int, d int) Datetime {
	tt := time.Time(*t)
	tt.AddDate(y, m, d)
	return Datetime(tt)
}