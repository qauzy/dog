package time

import (
	"database/sql/driver"
	"errors"
	"time"
)

type LocalDate time.Time

const (
	dateFormart = "2006-01-02"
)

func (t *LocalDate) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dateFormart+`"`, string(data), time.Local)
	*t = LocalDate(now)
	return
}

func (t LocalDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dateFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, dateFormart)
	b = append(b, '"')
	return b, nil
}

func (t LocalDate) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(dateFormart), nil
}

func (t *LocalDate) Scan(v interface{}) error {
	switch vt := v.(type) {
	case string:
		// 字符串转成 time.Time 类型
		tTime, _ := time.Parse(dateFormart, vt)
		*t = LocalDate(tTime)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}
