package utils

import (
	"strconv"
	"strings"
	"time"
)

func TimeToFormatString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func TimeFromFormatString(timeStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	return t, err
}

func IntArrayToString(array []int) string {
	arrayStr := make([]string, len(array))
	for i := 0; i < len(array); i++ {
		arrayStr[i] = strconv.Itoa(array[i])
	}
	return "[" + strings.Join(arrayStr, ",") + "]"
}

func StringArrayToString(array []string) string {
	return "[" + strings.Join(array, ",") + "]"
}
