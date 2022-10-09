package main

import "fmt"

// 创建一个数组类
type Arr []interface{}

func main() {
	// 创建空数组
	// 类型，len，cap
	arr1 := make(Arr, 0, 10)
	arr1.Info()
	arr1 = append(arr1, 1)
	arr1.Info()
	fmt.Println(arr1[0])
	// 由于当前位置数据未赋值，因此会报错
	// fmt.Println(arr1[1])
	// 创建默认初始化的数组
	arr2 := make(Arr, 10)
	arr2.Info()
	// 测试追加数据
	arr2 = append(arr2, 1)
	// 发现数组出现了自动扩容，cap变为两倍大小，但是数组长度以最后一个元素的位置为准
	arr2.Info()
	// 可以访问
	arr2[10] = 10
	// 访问出错
	// arr2[11] = 11
	fmt.Println(arr2)
	// 测试切片
	arr3 := make(Arr,0)
	for i:=0; i<10; i++{
		arr3 = append(arr3,i)
	}
	fmt.Println(arr3[0:3])
	// 测试切片修改
	arr4 := arr3[6:10]
	fmt.Println(arr4)
	arr4[0] = -1
	fmt.Println(arr3)
	fmt.Println(arr4)
}

// 展示信息方法
func (a Arr) Info(){
	fmt.Printf("len:%d, cap:%d\n", len(a), cap(a))
}
