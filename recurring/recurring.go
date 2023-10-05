package recurring

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"torch/torch-server/optional"
)

type Recurring struct {
	Times     optional.NullUint   `gorm:"column:rec_times" json:"times"`
	Period    optional.NullString `gorm:"column:rec_period" json:"period"`
	Progress  optional.NullUint   `gorm:"column:rec_progress" json:"progress"`
	UpdatedAt optional.NullString `gorm:"column:rec_updated_at" json:"updatedAt"`
}

type NullRecurring struct {
	Val     Recurring
	IsValid bool
}

func NewNullRecurring(val interface{}) NullRecurring {
	ni := NullRecurring{}
	ni.Set(val)
	return ni
}

func (ni *NullRecurring) Scan(value interface{}) error {
	fmt.Printf("\n Scan START : %v \n", value)

	if val, ok := value.(NullRecurring); ok {
		fmt.Printf("\n Scan: %v %v \n", val.Val, val.IsValid)
		ni.Val, ni.IsValid = val.Val, val.IsValid
	}

	return nil
}

func (ni NullRecurring) Value() (driver.Value, error) {
	fmt.Printf("\n VALUE : %v \n", ni)

	if !ni.IsValid {
		return nil, nil
	}
	return ni.Val, nil
}

func (ni *NullRecurring) Set(val interface{}) {
	fmt.Printf("\n SET : %v \n", ni)
	ni.Val, ni.IsValid = val.(NullRecurring).Val, val.(NullRecurring).IsValid
}

func (ni NullRecurring) MarshalJSON() ([]byte, error) {
	if !ni.IsValid {
		return []byte(`null`), nil
	}

	return json.Marshal(ni.Val)
}

func (ni *NullRecurring) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == `null` {
		ni.IsValid = false
		return nil
	}

	var recurring Recurring
	if err := json.Unmarshal(data, &recurring); err != nil {
		return err
	}

	ni.Val = recurring
	ni.IsValid = true

	return nil
}

func (ni NullRecurring) String() string {
	if !ni.IsValid {
		return `<nil>`
	}

	return fmt.Sprintf("Times: %v, Period: %v, Progress: %v", ni.Val.Times, ni.Val.Period, ni.Val.Progress)
}

// func (ni Recurring) MarshalJSON() ([]byte, error) {
// 	if ni.Times.Val == 0 {
// 		return []byte(`null`), nil
// 	}

// 	return json.Marshal(ni)
// }
