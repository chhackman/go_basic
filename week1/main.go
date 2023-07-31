package main

import (
	"fmt"
)

func main() {

	////删除切片特定下标元素的方法
	slice := []int{1, 2, 3, 4, 5, 6, 7}
	var sliceResult = DeleteElement(slice, 2)
	fmt.Printf("slice:%v\n", sliceResult)

	////使用泛型规定切片类型
	//sliceInt := []int{1, 2, 3, 4, 5, 6, 7}
	sliceFloat32 := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}

	var sliceResultT = DeleteElementT[float32](sliceFloat32, 2)
	fmt.Printf("sliceResultT:%v\n", sliceResultT)

	//泛型&&缩容，缩容策略len小于cap一半，cap改为原来cap的1/2
	var sliceRe = make([]int, 5, 10)
	i := 0
	for index, _ := range sliceRe {
		sliceRe[index] = i
		i++
	}
	fmt.Printf("sliceRe:%v, len:%d, cap:%d\n", sliceRe, len(sliceRe), cap(sliceRe))

	sliceResultTRe := DeleteElementTSub[int](sliceRe, 2)
	fmt.Printf("去掉元素后的列表sliceResultT:%v, len:%d, cap:%d\n", sliceResultTRe, len(sliceResultTRe), cap(sliceResultTRe))

}
