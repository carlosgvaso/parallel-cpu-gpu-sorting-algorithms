// Package quicksort provides a parallel quick sort implementation to sort
// integer arrays.
package quicksort

import (
	"math/rand"
	"sync"
)

// Sort sorts an array in place using the parallel quicksort algorithm.
//
// It takes an array as an input.
// It returns the input array sorted.
func Sort(arr []int) []int {
	var wg sync.WaitGroup // Wait group to synchronize parallel goroutines
	var n int = len(arr)  // Length of the array

	// Run quicksort
	wg.Add(1)
	quicksort(arr, 0, n-1, &wg)
	wg.Wait()

	return arr
}

// Quicksort is a regular quicksort implementation with random pivot and
// parallelized by goroutines at each recursive call
func quicksort(arr []int, p int, r int, wg *sync.WaitGroup) {
	defer wg.Done()

	if p < r {
		q := partition(arr, p, r)

		wg.Add(2)
		go quicksort(arr, p, q-1, wg)
		go quicksort(arr, q+1, r, wg)
	}
}

// Partition splits the input array using a randomized choice of a pivot.
func partition(arr []int, p int, r int) int {
	index := rand.Intn(r-p) + p
	pivot := arr[index]
	arr[index] = arr[r]
	arr[r] = pivot
	x := arr[r]
	j := p - 1
	i := p

	for i < r {
		if arr[i] <= x {
			j++

			tmp := arr[j]
			arr[j] = arr[i]
			arr[i] = tmp
		}

		i++
	}

	temp := arr[j+1]
	arr[j+1] = arr[r]
	arr[r] = temp

	return j + 1
}
