package prometheus

import (
	"fmt"
	"testing"
)

func removeDuplicates(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	l := 1
	for r := 1; r < len(nums); r++ {
		if nums[r] != nums[r-1] {
			nums[l] = nums[r]
			l++
		}
	}
	fmt.Println(nums)
	return l
}
func Test2(t *testing.T) {
	duplicates := removeDuplicates([]int{1, 1, 2})
	fmt.Println(duplicates)
}
