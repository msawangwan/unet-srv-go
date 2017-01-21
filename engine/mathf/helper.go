package mathf

// truncate a float to two levels of precision
func Truncate32(f float32) float32 {
	return float64(int(f*100)) / 100
}

// truncate a float to two levels of precision
func Truncate64(f float64) float64 {
	return float64(int(f*100)) / 100
}
