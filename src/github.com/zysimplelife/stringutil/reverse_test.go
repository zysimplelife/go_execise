package stringutil

import "testing"

func TestReverse(t *testing.T) {
	res := Reverse("Hello")
	if res != "olle" {
		t.Errorf("Failed to reverse hello, it got %s", res)
	}
}
