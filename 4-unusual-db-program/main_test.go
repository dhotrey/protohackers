package main

import (
	"fmt"
	"testing"
)

func TestParseKV(t *testing.T) {
	var tests = []struct {
		message                  string
		expectedKey, expectedVal string
	}{
		{"foo=bar", "foo", "bar"},
		{"foo=bar=baz", "foo", "bar=baz"},
		{"foo=", "foo", ""},
		{"foo===", "foo", "=="},
		{"=foo", "", "foo"},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("%d : %s", i, tt.message)
		t.Run(testname, func(t *testing.T) {
			parsedKey, parsedVal := parseKeyValue(tt.message)
			if parsedKey != tt.expectedKey && parsedVal != tt.expectedVal {
				t.Errorf("got key %s , val %s | want key %s , val %s", parsedKey, parsedVal, tt.expectedKey, tt.expectedVal)
			}
		})

	}

}
