package helpers

import (
	"encoding/json"
	"fmt"
)

func PrintJson(s interface{}) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
