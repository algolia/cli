package shared

import "math"

// Int32 safely converts an int to an int32, clamping to the int32 range so the
// conversion can never overflow. The Agent Studio SDK takes pagination and
// filter values as int32 while our flags are parsed as int; this keeps gosec
// (G115) satisfied without scattering //nolint directives across call sites.
func Int32(v int) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v)
}
