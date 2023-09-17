package optional

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

type NullUint struct {
	Val   uint
	Valid bool
}

func NewNullUint(val interface{}) NullUint {
	ni := NullUint{}
	ni.Set(val)
	return ni
}

func (ni *NullUint) Scan(value interface{}) error {
	data, ok := value.([]uint8)
	if !ok {
		return nil
	}

	val, err := strconv.ParseUint(string(data), 10, 64)
	if err == nil {
		ni.Val = uint(val)
		ni.Valid = true
	}
	return nil
}

func (ni NullUint) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return int64(ni.Val), nil
}

func (ni *NullUint) Set(val interface{}) {
	ni.Val, ni.Valid = val.(uint)
}

func (ni NullUint) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte(`null`), nil
	}

	return json.Marshal(ni.Val)
}

func (ni *NullUint) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.Valid = false
		return nil
	}

	val, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		ni.Valid = false
		return err
	}

	ni.Val = uint(val)
	ni.Valid = true

	return nil
}

func (ni NullUint) String() string {
	if !ni.Valid {
		return `<nil>`
	}

	return strconv.FormatUint(uint64(ni.Val), 10)
}
