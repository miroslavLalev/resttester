package stress_test

import (
	"reflect"
	"testing"
)

func AssertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Fatalf("%v should be nil", v)
	}
}

func AssertEqual(t *testing.T, v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Fatalf("Values %v and %v are not equal", v1, v2)
	}
}

func StringPtr(s string) *string {
	return &s
}
