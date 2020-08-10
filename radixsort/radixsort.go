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

// Sort sorts an array of positive integers in ascending order using the
// parallel most significant digit radix sort algorithm.
//
// arr is the positive integer input array to sort.
// k is the number of digits of the largest integer in the input array.
// It returns the input array sorted in ascending order.
func Sort(arr []int, k int) []int {
	var wg sync.WaitGroup // Wait group to synchronize parallel goroutines
	//var k int = 3         // Max number of digits in any array element (assumes max array entry is 999)
	/* I ended up not needing channels, but I'm keeping this just in case.
	// Channel needed to call radixsort()
	var returnChan chan []int
	*/

	// Run quicksort
	wg.Add(1)
	arr = radixsort(arr, 1, k, &wg) // , returnChan
	wg.Wait()

	/* I ended up not needing channels, but I'm keeping this just in case.
	// Receive the sorted array through the channel
	arr <- returnChan
	*/
	return arr
}

// Radixsort is a most significant digit radixsort implementation with
// parallelized by goroutines to fill the buckets and at each recursive call.
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
func radixsort(arr []int, l int, k int, wg *sync.WaitGroup) []int { //, ch chan []int
	defer wg.Done()

	// Check if we got just one element in the bucket
	if len(arr) == 1 {
		/* I ended up not needing channels, but I'm keeping this just in case.
		// Send the sorted array through the channel
		ch <- arr
		*/
		return arr
	}

	// Wait group to syncronize concurrent bucket operations using goroutines
	var wgBuckets sync.WaitGroup

	// Buckets is an array of slices. This data structure will be written
	// concurrently, so we need mutexes to make sure 2 or more goroutines don't
	// write to it at the same time. Since go allows to write concurrently to
	// different array entries, we create a mutex for each bucket.
	var buckets [numBuckets][]int
	var bucketLocks [numBuckets]sync.Mutex

	/* I ended up not needing channels, but I'm keeping this just in case.
	// Similarly, we create a channel for each bucket, so we can retrieve the
	// values of the recursive calls, which are called concurrently as
	// goroutines.
	var bucketChans [numBuckets]chan []int
	for i := range bucketChans {
		wgBuckets.Add(1)
		go func(i int) {
			defer wgBuckets.Done()

			bucketChans[i] = make(chan []int)
		}(i)
	}
	wgBuckets.Wait()
	*/

	// Place the elements in the buckets concurrently
	for _, v := range arr {
		wgBuckets.Add(1)
		go func(v int) {
			defer wgBuckets.Done()

			// Get the lth most significant digit, d, of the element. Elements with
			// less digits than the largest element are zero-padded.
			d := int((v / int(math.Pow10((k-1)/l))) % 10)

			// Place element in bucket buckets[d]
			bucketLocks[d].Lock()
			buckets[d] = append(buckets[d], v)
			bucketLocks[d].Unlock()
		}(v)
	}
	wgBuckets.Wait()

	if l <= k {
		// Concurrent recursive call
		for _, bucket := range buckets {
			// Only recurse if bucket is not empty
			if len(bucket) > 0 {
				wgBuckets.Add(1)
				// The bucket that we are passing is nothing but a slice that
				// the radixsort recursive call will sort. Since radixsort()
				// saves the sorted array back to the input array passed to it,
				// there is no need to return/get the sorted array, because it
				// will already be saved in the bucket.
				go radixsort(bucket, l+1, k, &wgBuckets) //, bucketChans[i]
			}
		}
		wgBuckets.Wait()

		/* I ended up not needing channels, but I'm keeping this just in case.
		// Get results. This assumes numBuckets = 10 (radix of 10).
		// Use channels to return data from radixsort()
		for i := 0; i < numBuckets; i++ {
			select {
			case arrTmp := <-bucketChans[0]:
				buckets[0] = arrTmp
			case arrTmp := <-bucketChans[1]:
				buckets[1] = arrTmp
			case arrTmp := <-bucketChans[2]:
				buckets[2] = arrTmp
			case arrTmp := <-bucketChans[3]:
				buckets[3] = arrTmp
			case arrTmp := <-bucketChans[4]:
				buckets[4] = arrTmp
			case arrTmp := <-bucketChans[5]:
				buckets[5] = arrTmp
			case arrTmp := <-bucketChans[6]:
				buckets[6] = arrTmp
			case arrTmp := <-bucketChans[7]:
				buckets[7] = arrTmp
			case arrTmp := <-bucketChans[8]:
				buckets[8] = arrTmp
			case arrTmp := <-bucketChans[9]:
				buckets[9] = arrTmp
			}
		}
		*/
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

	/* I ended up not needing channels, but I'm keeping this just in case.
	// Send the sorted array through the channel
	ch <- arr
	*/
	return arr
}
