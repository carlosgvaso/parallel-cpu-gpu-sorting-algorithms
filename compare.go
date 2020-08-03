// Compare all included parallel sorting algorithms' performance.
package main

import (
	"fmt"

	"github.com/carlosgvaso/parallel-sort/bricksort"
)

// Main runs the tests to compare the sorting algorithms.
func main() {
	arrIn := []int{1, 2, 3, 4, 5}
	arrOut := make([]int, 5)

	arrOut = bricksort.Sort(arrIn)

	fmt.Printf("Brick sort:\tIn=%v\tOut=%v\n", arrIn, arrOut)
}
