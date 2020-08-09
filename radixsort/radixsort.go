// Package radixsort provides a parallel radix sort implementation to sort
// arrays of positive integers in ascending order.
package radixsort

import (
	"math"
	"sync"
)

// NumBuckets is the number of buckets.
//
// Since we are sorting positive integers (e.i decimal numbers), it is 10.
const numBuckets int = 10

// Sort sorts an array in ascending order using the parallel most significant
// digit radix sort algorithm.
//
// This function assumes the maximum integer in the array to sort has 3 digits.
// Therefore, it can only sort arrays with integer entries in the range
// [0, 999].
//
// arr is the input array to sort.
// It returns the input array sorted in ascending order.
func Sort(arr []int) []int {
	var wg sync.WaitGroup // Wait group to synchronize parallel goroutines
	var k int = 3         // Max number of digits in any array element (assumes max array entry is 999)

	// Run quicksort
	wg.Add(1)
	radixsort(arr, 1, k, &wg)
	wg.Wait()

	return arr
}

// Radixsort is a most significant digit radixsort implementation with
// parallelized by goroutines at each recursive call.
//
// It is based in the algorithm presented in Alg. 2 of PARADIS: A PARALLEL
// IN-PLACE RADIX SORT ALGORITHM by Rolland He: PARADIS:
// https://stanford.edu/~rezab/classes/cme323/S16/projects_reports/he.pdf
//
// arr is the input array to sort.
// l is the current most significant digit, where l=1 is the most significant
// digit of the largest integer in the array, and l=k-1 is the least significant
// digit.
// k is the maximum number of digits in the array.
// wg is a sync.WaitGroup for synchronization of the goroutines.
// It returns the array sorted in ascending order.
func radixsort(arr []int, l int, k int, wg *sync.WaitGroup) []int {
	defer wg.Done()

	// Check if we got just one element in the bucket
	if len(arr) == 1 {
		return arr
	}

	// Buckets is an array of slices
	var buckets [numBuckets][]int

	// Place the elements in the buckets
	for _, v := range arr {
		// Get the lth most significant digit, d, of the element. Elements with
		// less digits than the largest element are zero-padded.
		d := int((v / int(math.Pow10((k-1)/l))) % 10)

		// Place element in bucket buckets[d]
		buckets[d] = append(buckets[d], v)
	}

	if l <= k {
		for i, bucket := range buckets {
			wg.Add(1)
			buckets[i] = radixsort(bucket, l+1, k, wg)
		}
	}

	// Replace arr elements with elements from buckets in the same order
	var i int = 0
	// Iterate over all the buckets in ascendig order
	for _, bucket := range buckets {
		// Iterate over all the entries in a bucket in ascending order
		for _, v := range bucket {
			// Add entries to arr
			arr[i] = v
			i++
		}
	}

	return arr
}
