package bitonicsort

import (
	"sync"
)

//Adding boolean for Ascending and descending order
const (
	ASC  bool = true
	DESC bool = false
)

// Sort sorts an array in place using the parallel bitonic sort algorithm.
func Sort(arr []int, diff int) []int {
	orderby := true
	bitonicSort(arr, orderby)
	arr = arr[diff:]
	return arr
}

// CheckAndAppendZeros updates the array by making its length exponential of 2 by appending 0's
// returns the array and difference of updated and acutal size
func CheckAndAppendZeros(arr []int) ([]int, int) {
	i := 2
	size := len(arr)
	for size > i {
		i *= 2
	}
	i = i - size
	var a = make([]int, i)
	arr = append(arr, a...)
	return arr, i
}

// bitonicSort  will return the sorted array based on the input array
func bitonicSort(arr []int, orderby bool) {
	if len(arr) < 2 {
		return
	}

	middle := len(arr) / 2
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bitonicSort(arr[:middle], ASC)
	}()

	go func() {
		defer wg.Done()
		bitonicSort(arr[middle:], DESC)
	}()
	wg.Wait()
	bitonicMerge(arr, orderby)
}

func bitonicCompare(arr []int, orderby bool) {
	middle := len(arr) / 2
	for i := 0; i < middle; i++ {
		if (arr[i] > arr[i+middle]) == orderby {
			arr[i], arr[i+middle] = arr[i+middle], arr[i]
		}
	}
}

func bitonicMerge(arr []int, orderby bool) {
	bitonicCompare(arr, orderby)
	middle := len(arr) / 2
	if middle > 1 {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			bitonicMerge(arr[:middle], orderby)
		}()
		go func() {
			defer wg.Done()
			bitonicMerge(arr[middle:], orderby)
		}()
		wg.Wait()

	}
}
