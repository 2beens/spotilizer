package services_test

import "testing"

func TestDummy(t *testing.T) {
	t.Log(" > starting test dummy ...")
	s1 := ""
	if len(s1) > 0 {
		t.Errorf(" NO!")
	}
}
