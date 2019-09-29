package bbrpc

import (
	"encoding/json"
)

func toJSONIndent(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// pointer of string
func ps(s string) *string {
	return &s
}

// pointer of bool
func pbool(b bool) *bool {
	return &b
}
