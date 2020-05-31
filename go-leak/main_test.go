package main

import (
	"testing"

	"go.uber.org/goleak"
)

// TestA ...
func TestA(t *testing.T) {
	defer goleak.VerifyNone(t)

	Run()
	// test logic here.
}
