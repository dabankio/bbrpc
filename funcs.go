package bbrpc

import (
	"encoding/json"
)

func toJSONIndent(v interface{}) string { b, _ := json.MarshalIndent(v, "", "  "); return string(b) }

func ps(s string) *string      { return &s }
func pbool(b bool) *bool       { return &b }
func puint(i uint) *uint       { return &i }
func pstring(s string) *string { return &s }
func pint(i int) *int          { return &i }
