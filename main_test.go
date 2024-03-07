package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var s *string
	fmt.Println(s == nil)
	fmt.Println(*s)
}
