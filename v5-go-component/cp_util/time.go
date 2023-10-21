package cp_util

import (
	"errors"
	"fmt"
	timeconv "github.com/Andrew-M-C/go.timeconv"
	"time"
)

var timeFormat = "2006-01-02 15:04:05"

//获取当前时间，格式：yyyy-MM-dd HH:mm:ss
func GetTimeByString() string {
	t := time.Now()
	return FormatTimeToString(t)
}

/*
  将日期字符串解析为系统time
  author: max.feng
  date: 2016-12-02
*/
func ParseStringToTime(str string) (time.Time, error) {
	return ParseStringToTimeUseFormat(str, "")
}

/*
  两日期相差多少天
  author: csdn up
  date: 2021-8-14
*/
func GetDateDiff(startTime, endTime time.Time) (int64, error) {
	timeLayout := "2006-01-02  15:04:05"
	//// 转成时间戳
	startUnix, err := time.ParseInLocation(timeLayout, FormatTimeToString(startTime), time.Local)
	if err != nil{
		return 0,errors.New("start时间格式错误")
	}
	//
	endUnix, err := time.ParseInLocation(timeLayout, FormatTimeToString(endTime), time.Local)
	if err != nil{
		return 0,errors.New("end时间格式错误")
	}

	// 求相差天数
	date :=	(endUnix.Unix() - startUnix.Unix()) / 86400
	return date, nil
}

/*
  将日期字符串用自定义时间格式解析为系统time
  author: max.feng
  date: 2016-12-02
*/
func ParseStringToTimeUseFormat(str, format string) (time.Time, error) {
	var err error
	if format == "" {
		format = timeFormat
	}
	if loc, err := time.LoadLocation("Local"); err == nil {
		if t, err := time.ParseInLocation(format, str, loc); err == nil {
			return t, nil
		}
	}
	return time.Now(), err
}

//将指定时间转化为指定格式字符串：yyyy-MM-dd HH:mm:ss
func FormatTimeToString(t time.Time) string {
	return t.Format(timeFormat)
}

//避免原生包AddDate的坑 原生的AddDate会出现2022-01-31加上1个月后，会返回3月3日的现象
func AddDate(t time.Time, years int, months int, days int) time.Time {
	return timeconv.AddDate(t, years, months, days)
}

func ListYearMonth(from, to int64, maxday int) ([]string, error) {
	var newStr string

	fromTime := time.Unix(from, 0)
	toTime := time.Unix(to, 0)

	if from > time.Now().Unix() {
		return nil, errors.New("起始时间不能大于当前时间")
	} else if fromTime.Sub(toTime) > 0 {
		return nil, errors.New("起始时间不能大于结束时间")
	} else if to > time.Now().Unix(){
		to = time.Now().Unix()
	}

	if fromTime.AddDate(0, 0, maxday).Before(toTime) {
		return nil, errors.New(fmt.Sprintf("最多只能查询%d天内的记录", maxday))
	}

	fromStr := fmt.Sprintf("%d_%d", fromTime.Year(), fromTime.Month())
	endStr := fmt.Sprintf("%d_%d", toTime.Year(), toTime.Month())

	yearMonthList := make([]string, 0)
	yearMonthList = append(yearMonthList, fromStr)

	newPoint := fromTime

	for {
		newPoint = AddDate(newPoint, 0, 1, 0)

		if newPoint.Sub(toTime) >= 0 {
			if endStr != newStr && endStr != fromStr {
				yearMonthList = append(yearMonthList, endStr)
			}
			break
		}

		newStr = fmt.Sprintf("%d_%d", newPoint.Year(), newPoint.Month())
		yearMonthList = append(yearMonthList, newStr)
	}

	if len(yearMonthList) == 0 {
		return nil, errors.New("时间格式错误")
	}

	return yearMonthList, nil
}