package array

import "errors"

// RemoveItemFromArray 移除数组中指定索引的元素，并返回移除后的新数组。
func RemoveItemFromArray[T any](arr []T, index int) []T {
	if index < 0 || index >= len(arr) {
		return arr
	}
	newArr := make([]T, len(arr)-1)
	copy(newArr, arr[:index])
	copy(newArr[index:], arr[index+1:])
	return newArr
}

func InsertIntoArray[T any](arr []T, index int, value T) ([]T, error) {
	// 检查索引是否有效
	if index < 0 || index > len(arr) {
		return nil, errors.New("range index error")
	}

	// 创建一个新的数组，长度为原数组长度加一
	newLength := len(arr) + 1
	newArray := make([]T, newLength)

	// 复制插入点之前的元素
	copy(newArray, arr[:index])

	// 插入新元素
	newArray[index] = value

	// 复制插入点之后的元素
	copy(newArray[index+1:], arr[index:])

	return newArray, nil
}
