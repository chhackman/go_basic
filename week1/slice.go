package main

func DeleteElement(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice
	}

	// 将索引之后的元素向前移动一个位置覆盖掉要删除的元素
	copy(slice[index:], slice[index+1:])

	// 返回删除元素后的新切片
	return slice[:len(slice)-1]
	//println(111)
}

// 通用的删除函数，支持Number类型的切片，Number可以定义为接口
func DeleteElementT[T Number](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}

	// 将索引之后的元素向前移动一个位置覆盖掉要删除的元素
	copy(slice[index:], slice[index+1:])

	// 返回删除元素后的新切片
	return slice[:len(slice)-1]
}

type Number interface {
	~int | int64 | float64 | float32 | int32 | byte
}

// 通用的删除函数，支持任意类型的切片
func DeleteElementTSub[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}

	// 将索引之后的元素向前移动一个位置覆盖掉要删除的元素
	copy(slice[index:], slice[index+1:])

	// 返回删除元素后的新切片
	result := slice[:len(slice)-1]

	// 如果切片的长度明显小于容量的一半，进行缩容
	if len(result) < cap(result)/2 {
		newSlice := make([]T, len(result), cap(result)/2)
		copy(newSlice, result)
		result = newSlice
	}

	return result
}
