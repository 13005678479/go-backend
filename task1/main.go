package main

import "fmt"

// twoSum 接收整数数组和目标值，返回和为目标值的两个整数的下标
func twoSum(nums []int, target int) []int {
	// 创建一个哈希表用于存储数值到下标的映射
	numMap := make(map[int]int)
	// 遍历数组
	for i, num := range nums {
		// 计算当前数值需要的互补数
		complement := target - num
		// 检查互补数是否在哈希表中
		if j, exists := numMap[complement]; exists {
			return []int{j, i}
		}
		// 如果不存在，将当前数值和下标存入哈希表
		numMap[num] = i
	}
	return []int{}
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(nums, target)
	fmt.Println(result)
}
