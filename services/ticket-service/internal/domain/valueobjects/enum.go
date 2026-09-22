package valueobjects

import "slices"

func isOneOf[T ~string](value T, allowed ...T) bool {
	return slices.Contains(allowed, value)
}
