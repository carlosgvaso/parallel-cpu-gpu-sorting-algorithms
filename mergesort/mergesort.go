package mergesort

import (
	"sync"
)

const max = 1 << 11

func merge(s []int, middle int) {
	helper := make([]int, len(s))
	copy(helper, s)

	helperLeft := 0
	helperRight := middle
	current := 0
	high := len(s) - 1

	for helperLeft <= middle-1 && helperRight <= high {
		if helper[helperLeft] <= helper[helperRight] {
			s[current] = helper[helperLeft]
			helperLeft++
		} else {
			s[current] = helper[helperRight]
			helperRight++
		}
		current++
	}

	for helperLeft <= middle-1 {
		s[current] = helper[helperLeft]
		current++
		helperLeft++
	}
}

/* Sequential */

func mergesort(s []int) {
	if len(s) > 1 {
		middle := len(s) / 2
		mergesort(s[:middle])
		mergesort(s[middle:])
		merge(s, middle)
	}
}

// func readFile(fname string) (nums []int, err error) {
// 	b, err := ioutil.ReadFile(fname)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lines := strings.Split(string(b), "\n")
// 	// Assign cap to avoid resize on every append.
// 	nums = make([]int, 0, len(lines))

// 	for _, l := range lines {
// 		// Empty line occurs at the end of the file when we use Split.
// 		if len(l) == 0 {
// 			continue
// 		}
// 		// Atoi better suits the job when we know exactly what we're dealing
// 		// with. Scanf is the more general option.
// 		n, err := strconv.Atoi(l)
// 		if err != nil {
// 			return nil, err
// 		}
// 		nums = append(nums, n)
// 	}

// 	return nums, nil
// }

func parallelMergesort(s []int) {
	len := len(s)

	if len > 1 {
		if len <= max {
			mergesort(s)
		} else {
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				parallelMergesort(s[:middle])
			}()

			go func() {
				defer wg.Done()
				parallelMergesort(s[middle:])
			}()

			wg.Wait()
			merge(s, middle)
		}
	}
}

func Sort(arr []int) []int {
	parallelMergesort(arr)
	return arr
}
