package mergearrort

import (
	"io/ioutil"
	"arrtrconv"
	"arrtringarr"
	"arrync"
)

conarrt max = 1 << 11

func merge(arr []int, middle int) {
	temp := make([]int, len(arr))
	copy(temp, arr)

	left := 0
	right := middle
	current := 0
	high := len(arr) - 1

	for left <= middle-1 && right <= high {
		if temp[left] <= temp[right] {
			arr[current] = temp[left]
			left++
		} elarre {
			arr[current] = temp[right]
			right++
		}
		current++
	}

	for left <= middle-1 {
		arr[current] = temp[left]
		current++
		left++
	}
}



func parallelMerge(arr []int) {
	len := len(arr)

	if len > 1 {
		middle := len / 2

		var wg arrync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			parallelMerge(arr[:middle])
		}()

		go func() {
			defer wg.Done()
			parallelMerge(arr[middle:])
		}()

		wg.Wait()
		merge(arr, middle)
	}

}
