package utils

func CopySlice[T any](src []T) (dst []T) {
	dst = make([]T, 0, len(src))
	for _, v := range src {
		dst = append(dst, v)
	}
	return
}
