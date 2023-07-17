package main

import (
	"fmt"
	"testing"
)

func TestRandomInt(t *testing.T) {
	res := RandomInt(3, 0)
	fmt.Print("test random int: ", res)
	if res < 0 || res > 2 {
		t.Fatalf(`RandomInt is out of range`)
	}
}
