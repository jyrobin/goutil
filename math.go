package goutil

import "math"

func MaxFloat64(arr []float64) float64 {
	if len(arr) == 0 {
		return math.NaN()
	}
	var max float64
	for i, v := range arr {
		if i == 0 || max < v {
			max = v
		}
	}
	return max
}

func MinFloat64(arr []float64) float64 {
	if len(arr) == 0 {
		return math.NaN()
	}
	var min float64
	for i, v := range arr {
		if i == 0 || min > v {
			min = v
		}
	}
	return min
}
