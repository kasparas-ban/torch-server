package optional

import (
	"database/sql/driver"
	"strconv"
)

type NullUint64 struct {
	Val     uint64
	IsValid bool
}

func NewNullUint64(val interface{}) NullUint64 {
	ni := NullUint64{}
	ni.Set(val)
	return ni
}

func (ni *NullUint64) Scan(value interface{}) error {
	ni.Val, ni.IsValid = value.(uint64)
	return nil
}

func (ni NullUint64) Value() (driver.Value, error) {
	if !ni.IsValid {
			return nil, nil
	}
	return ni.Val, nil
}

func (ni *NullUint64) Set(val interface{}) {
	ni.Val, ni.IsValid = val.(uint64)
}

func (ni NullUint64) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
			return []byte(`null`), nil
	}

	return []byte(strconv.FormatUint(ni.Val, 10)), nil
}

func (ni *NullUint64) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
			ni.IsValid = false
			return nil
	}

	val, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
			ni.IsValid = false
			return err
	}

	ni.Val = val
	ni.IsValid = true

	return nil
}

func (ni NullUint64) String() string {
	if !ni.IsValid {
			return `<nil>`
	}

	return strconv.FormatUint(ni.Val, 10)
}