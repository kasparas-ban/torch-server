package optional

import (
	"database/sql/driver"
	"strconv"
)

// type NullUint64 struct {
// 	Val     uint64
// 	IsValid bool
// }

// func NewNullUint64(val interface{}) NullUint64 {
// 	ni := NullUint64{}
// 	ni.Set(val)
// 	return ni
// }

// func (ni *NullUint64) Scan(value interface{}) error {
// 	data, ok := value.([]uint8)
// 	if !ok {
// 		return nil
// 	}

// 	val, err := strconv.ParseUint(string(data), 10, 64)
// 	if err == nil {
// 		ni.Val = val
// 		ni.IsValid = true
// 	}
// 	return nil
// }

// func (ni NullUint64) Value() (driver.Value, error) {
// 	if !ni.IsValid {
// 		return nil, nil
// 	}
// 	return ni.Val, nil
// }

// func (ni *NullUint64) Set(val interface{}) {
// 	ni.Val, ni.IsValid = val.(uint64)
// }

// func (ni NullUint64) MarshalJSON() ([]byte, error) {
// 	if !ni.IsValid {
// 		return []byte(`null`), nil
// 	}

// 	return json.Marshal(ni.Val)
// }

// func (ni *NullUint64) UnmarshalJSON(data []byte) error {
// 	if data == nil || string(data) == `null` {
// 		ni.IsValid = false
// 		return nil
// 	}

// 	val, err := strconv.ParseUint(string(data), 10, 64)
// 	if err != nil {
// 		ni.IsValid = false
// 		return err
// 	}

// 	ni.Val = val
// 	ni.IsValid = true

// 	return nil
// }

// func (ni NullUint64) String() string {
// 	if !ni.IsValid {
// 		return `<nil>`
// 	}

// 	return strconv.FormatUint(ni.Val, 10)
// }

type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

func (n *NullUint64) Scan(value any) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		n.Uint64 = 0
		n.Valid = false
		return nil
	}

	intValue, err := strconv.ParseUint(string(bytes), 10, 64)
	if err != nil {
		return nil
	}
	n.Uint64 = intValue
	n.Valid = false

	return nil
}

func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint64, nil
}
