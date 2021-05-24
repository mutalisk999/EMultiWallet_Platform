package coin

import (
	"fmt"
	"testing"
)

func TestFromPrecisionAmount(t *testing.T) {
	value := FromPrecisionAmount(3230000000, 8)
	fmt.Println(value)
}

func TestToPrecisionAmount(t *testing.T) {
	amount, _ := ToPrecisionAmount("1.234", 8)
	fmt.Println(amount)
}
