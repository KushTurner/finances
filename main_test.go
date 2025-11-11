package main

import "testing"

func TestAdd(t *testing.T) {
	expected := 8
	actual := Add(5, 3)

	if expected != actual {
		t.Errorf("Result was incorrect, got: %d, want: %d.", actual, expected)
	}
}
