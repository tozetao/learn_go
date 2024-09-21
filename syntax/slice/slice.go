package slice

import (
	"errors"
)

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Add[T Number](slice []T, idx int, val T) ([]T, error) {
	if idx < 0 || idx > len(slice) {
		return nil, errors.New("the index out of range")
	}

	result := make([]T, 0, cap(slice)+1)

	for i := 0; i < idx; i++ {
		result = append(result, slice[i])
	}
	result = append(result, val)

	for i := idx; i < len(slice); i++ {
		result = append(result, slice[i])
	}
	return result, nil
}

func DeleteAt[T any](slice []T, idx int) ([]T, error) {
	if idx < 0 || idx >= len(slice) {
		return nil, errors.New("the index out of range")
	}

	len := len(slice)

	for i := idx; i+1 < len; i++ {
		slice[i] = slice[i+1]
	}

	return slice[:len-1], nil
}

// func Delete(slice []int, idx int) ([]int, error) {
// 	if idx < 0 || idx >= len(slice) {
// 		return nil, errors.New("the idx out of range")
// 	}

// 	len := len(slice)
// 	if len == 0 || len == 1 {
// 		return slice[0:0], nil
// 	}

// 	// 截取idx下标之前的元素
// 	result := slice[0:idx]

// 	if idx == len-1 {
// 		return result, nil
// 	}

// 	result = append(result, slice[idx+1:]...)
// 	return result, nil
// }
