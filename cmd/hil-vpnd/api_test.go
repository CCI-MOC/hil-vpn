package main

import (
	"testing"
)

func TestCreate(t *testing.T) {
	ops := NewMockPrivOps()
	handler := makeHandler(ops)

	// Avoid compiler error re: unused variable; we'll do something
	// with it later:
	_ = handler
}
