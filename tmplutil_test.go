package main

import (
	"testing"
)

func TestName(t *testing.T) {
	if !HasField(struct {
		Name  string
		Grade int
	}{
		Name:  "yaozongyou",
		Grade: 100,
	}, "Name") {
		t.Error("Response from HasField is unexpected value")
	}

	if str := SprintField(struct {
		Name  string
		Grade int
	}{
		Name:  "yaozongyou",
		Grade: 100,
	}, "Name"); str != "yaozongyou" {
		t.Error("Response from SprintField is unexpected value")
	}
}
