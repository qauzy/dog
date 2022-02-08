package time

import (
	"database/sql/driver"
	"errors"
	"time"
)

type LocalDateTime time.Time

const (
	dateTimeFormart = "2006-01-02"
)

func (t *LocalDateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dateTimeFormart+`"`, string(data), time.Local)
	*t = LocalDateTime(now)
	return
}

func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dateTimeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, dateTimeFormart)
	b = append(b, '"')
	return b, nil
}

func (t LocalDateTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(dateTimeFormart), nil
}

func (t *LocalDateTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case string:
		// 字符串转成 time.Time 类型
		tTime, _ := time.Parse(dateTimeFormart, vt)
		*t = LocalDateTime(tTime)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}
