package bbrpc

import (
	"reflect"
	"runtime/debug"
	"testing"
)

func tShouldNil(t *testing.T, v interface{}, args ...interface{}) {
	if v != nil {
		debug.PrintStack()
		t.Fatalf("[test assert] should nil, but got: %v, %v", v, args)
	}
}

func tShouldTrue(t *testing.T, b bool, args ...interface{}) {
	if !b {
		debug.PrintStack()
		t.Fatalf("[test assert] should true, args: %v", args)
	}
}

func tShouldNotZero(t *testing.T, v interface{}, args ...interface{}) {
	value := reflect.ValueOf(v)
	if value.IsZero() {
		debug.PrintStack()
		t.Fatalf("[test assert] should not [zero value], %v", args)
	}
}
