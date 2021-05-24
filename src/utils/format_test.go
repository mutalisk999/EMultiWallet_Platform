package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeToFormatString(t *testing.T) {
	fmt.Println("timeStr:", TimeToFormatString(time.Now()))
}

func TestTimeFromFormatString(t *testing.T) {
	tm, err := TimeFromFormatString("2018-09-06 17:26:52")
	fmt.Println("tm:", tm)
	fmt.Println("err:", err)
}

func TestIntArrayToString(t *testing.T) {
	array := []int{1, 2, 3, 4}
	fmt.Println(IntArrayToString(array))
}

func TestStringArrayToString(t *testing.T) {
	array := []string{"1", "2", "3", "4"}
	fmt.Println(StringArrayToString(array))
}
