package util

import "encoding/json"

func WriteValueAsString(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}
