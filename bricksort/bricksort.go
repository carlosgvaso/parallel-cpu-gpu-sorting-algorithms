// Package bricksort provides a parallel brick sort implementation to sort
// integer arrays.
package bricksort

import (
	"sync"
)

// Sort sorts an array in place using the parallel brick sort algorithm.
//
// It takes an array as an input.
// It returns the input array sorted.
func Sort(arr []int) []int {
	var waitGroup sync.WaitGroup // Wait group to synchronize parallel goroutines
	var isSorted bool = false    // True if there weren't any swaps for an iter
	var n int = len(arr)         // Length of the array
	var iter int = 0             // Count the iterations it takes to get sorted

	// Iterate until the array is sorted (max n iters)
	for isSorted == false {
		iter++
		isSorted = true // Assume is sorted, and set false if there is a swap

		for i := 1; i < n; i += 2 {
			waitGroup.Add(1)
			go swap(arr, i, &isSorted, &waitGroup)
		}
		waitGroup.Wait()

		for i := 2; i < n; i += 2 {
			waitGroup.Add(1)
			go swap(arr, i, &isSorted, &waitGroup)
		}
		waitGroup.Wait()
	}

	return arr
}

// Swap checks if arr[i-1] > arr[i], and swaps their values and sets isSorted to
// false if true. It does nothing otherwise.
func swap(arr []int, i int, isSorted *bool, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	if arr[i-1] > arr[i] {
		tmp := arr[i]
		arr[i] = arr[i-1]
		arr[i-1] = tmp
		*isSorted = false
	}
}
