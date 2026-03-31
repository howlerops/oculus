package array

// Unique returns a deduplicated slice preserving order.
func Unique[T comparable](items []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// Filter returns items matching the predicate.
func Filter[T any](items []T, pred func(T) bool) []T {
	var result []T
	for _, item := range items {
		if pred(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms each item.
func Map[T any, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// Find returns the first matching item.
func Find[T any](items []T, pred func(T) bool) (T, bool) {
	for _, item := range items {
		if pred(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Contains checks if an item is in the slice.
func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

// Count returns how many items match the predicate.
func Count[T any](items []T, pred func(T) bool) int {
	n := 0
	for _, item := range items {
		if pred(item) {
			n++
		}
	}
	return n
}

// Last returns the last item or zero value.
func Last[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	return items[len(items)-1], true
}

// Chunk splits a slice into chunks of size n.
func Chunk[T any](items []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}
