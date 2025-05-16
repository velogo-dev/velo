package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v any) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(json))
}
