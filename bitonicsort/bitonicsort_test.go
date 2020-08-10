// Test parallel bitonic sort implementation
package bitonicsort

import (
	"reflect"
	"testing"
)

// TestSort checks Sort with a multitude of input arrays.
func TestSort(t *testing.T) {
	cases := []struct {
		in, want []int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7}, []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{[]int{7, 6, 5, 4, 3, 9, 2, 1, 0, 8}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[]int{1, 3, 5, 7, 6, 4, 2, 0}, []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{[]int{3, 0, 5, 7, 1, 6, 2, 4}, []int{0, 1, 2, 3, 4, 5, 6, 7}},
	}

	for _, c := range cases {
		var diff int
		arrIn := make([]int, len(c.in))

		// Copy and zero-pad (to make power of 2) input array to arrIn
		copy(arrIn, c.in)
		arrIn, diff = CheckAndAppendZeros(arrIn)

		got := Sort(arrIn, diff)
		want := c.want

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Sort (%v) == %v, want %v", c.in, got, want)
		}
	}
}
