package utils

import (
	"fmt"
	"testing"
)

func TestMakeError(t *testing.T) {
	InitGlobalError()
	err := MakeError(100000, "test", 1, 2, 3)
	fmt.Println(*err)
}
