package main

import (
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	t.Run("read non-existent file", func(t *testing.T) {
		_, err := readJSON("non-existent.json")
		if err == nil {
			t.Fatalf("readJSON: expected error, got nil")
		}
		if !strings.HasPrefix(err.Error(), "open") {
			t.Fatalf("readJSON: expected error to start with 'open', got %v", err)
		}
	})

	t.Run("read invalid file", func(t *testing.T) {
		_, err := readJSON("../../testdata/test_invalid.json")
		if err == nil {
			t.Fatalf("readJSON: expected error, got nil")
		}
		if !strings.HasPrefix(err.Error(), "unmarshal") {
			t.Fatalf("readJSON: expected error to start with 'unmarshal', got %v", err)
		}
	})

	t.Run("read valid file", func(t *testing.T) {
		actual, err := readJSON("../../testdata/test_int.json")
		if err != nil {
			t.Fatalf("readJSON: %v", err)
		}
		if len(actual) == 0 {
			t.Fatalf("readJSON: expected non-empty map, got empty map")
		}
		if _, ok := actual["IntField"]; !ok {
			t.Fatalf("readJSON: expected key IntField, got %v", actual)
		}
	})
}
