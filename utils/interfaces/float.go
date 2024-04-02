package interfaces

import (
	"errors"
	"math"
	"strconv"
)

// Float64 convert to float64
func Float64(v interface{}) (float64, error) {
	switch i := v.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		return math.NaN(), errors.New("math: square root of negative number")
	}
}

// Float32 convert to float32
func Float32(v interface{}) (float32, error) {
	switch i := v.(type) {
	case float64:
		return float32(i), nil
	case float32:
		return i, nil
	case uint:
		return float32(i), nil
	case int:
		return float32(i), nil
	case uint32:
		return float32(i), nil
	case int32:
		return float32(i), nil
	case uint64:
		return float32(i), nil
	case int64:
		return float32(i), nil
	case string:
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0, err
		}
		return float32(f), nil
	default:
		return 0, errors.New("math: square root of negative number")
	}
}
