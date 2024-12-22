package array

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
