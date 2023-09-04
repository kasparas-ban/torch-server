package optional

import (
	"database/sql/driver"
)

type NullString struct {
	Val     string
	IsValid bool
}

func NewNullString(val interface{}) NullString {
	ni := NullString{}
	ni.Set(val)
	return ni
}

func (ni *NullString) Scan(value interface{}) error {
	ni.Val, ni.IsValid = value.(string)
	return nil
}

func (ni NullString) Value() (driver.Value, error) {
	if !ni.IsValid {
		return nil, nil
	}
	return ni.Val, nil
}

func (ni *NullString) Set(val interface{}) {
	ni.Val, ni.IsValid = val.(string)
}

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
		return []byte(`null`), nil
	}

	return []byte(ni.Val), nil
}

func (ni *NullString) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.IsValid = false
		return nil
	}

	ni.Val = string(data)
	ni.IsValid = true

	return nil
}

func (ni NullString) String() string {
	if !ni.IsValid {
		return `<nil>`
	}

	return ni.Val
}
