package util

import (
	"fmt"
)

// MaxAndMin return the max and min number
func MaxAndMin(data []float64) (max, min float64) {
	if len(data) > 0 {
		max, min = data[0], data[0]
	}

	for _, item := range data {
		if item < min {
			min = item
		} else if item > max {
			max = item
		}
	}
	return
}

// PrintCollectTrend print the trend of data
func PrintCollectTrend(data []float64) (buf string) {
	max, min := MaxAndMin(data)

	unit := (max - min) / 100
	for _, num := range data {
		total := (int)(num / unit)
		if total == 0 {
			total = 1
		}
		arr := make([]int, total)
		for range arr {
			buf = fmt.Sprintf("%s*", buf)
		}
		buf = fmt.Sprintf("%s %.0f\n", buf, num)
	}
	return
}
