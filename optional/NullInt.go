package optional

import (
	"database/sql/driver"
	"strconv"
)

type NullInt struct {
	Val     int
	IsValid bool
}

func NewNullInt(val interface{}) NullInt {
	ni := NullInt{}
	ni.Set(val)
	return ni
}

func (ni *NullInt) Scan(value interface{}) error {
	ni.Val, ni.IsValid = value.(int)
	return nil
}

func (ni NullInt) Value() (driver.Value, error) {
	if !ni.IsValid {
			return nil, nil
	}
	return ni.Val, nil
}

func (ni *NullInt) Set(val interface{}) {
	ni.Val, ni.IsValid = val.(int)
}

func (ni NullInt) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
			return []byte(`null`), nil
	}

	return []byte(strconv.FormatInt(int64(ni.Val), 10)), nil
}

func (ni *NullInt) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
			ni.IsValid = false
			return nil
	}

	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
			ni.IsValid = false
			return err
	}

	ni.Val = int(val)
	ni.IsValid = true

	return nil
}

func (ni NullInt) String() string {
	if !ni.IsValid {
			return `<nil>`
	}

	return strconv.FormatInt(int64(ni.Val), 10)
}