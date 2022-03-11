package util

// This is used instead of `math.Min/Max` because of the
// unnecessary coercion from/to float64.
func Clamp(value, min, max int) int {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}
