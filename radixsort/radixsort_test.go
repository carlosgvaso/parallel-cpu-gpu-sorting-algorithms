// Test parallel radixsort implementation
package radixsort

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
		{[]int{7, 6, 5, 4, 3, 2, 1, 0}, []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{[]int{1, 3, 5, 7, 6, 4, 2, 0}, []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{[]int{3, 0, 5, 7, 1, 6, 2, 4}, []int{0, 1, 2, 3, 4, 5, 6, 7}},
		{[]int{999, 4, 295, 666, 43, 66, 6, 576}, []int{4, 6, 43, 66, 295, 576, 666, 999}},
	}

	for _, c := range cases {
		arrIn := make([]int, len(c.in))
		copy(arrIn, c.in)

		got := Sort(arrIn)
		want := c.want

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Sort (%v) == %v, want %v", c.in, got, want)
		}
	}
}
