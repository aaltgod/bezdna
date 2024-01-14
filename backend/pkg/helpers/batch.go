package helpers

func Batch[T any](values []T, size int) [][]T {
	if len(values) == 0 || size == 0 {
		return nil
	}

	result := make([][]T, 0, int(float64(len(values)/size)))

	for leftBoundary := 0; leftBoundary < len(values); leftBoundary += size {
		rightBoundary := leftBoundary + size
		if rightBoundary > len(values) {
			rightBoundary = len(values)
		}

		result = append(result, values[leftBoundary:rightBoundary])
	}

	return result
}
