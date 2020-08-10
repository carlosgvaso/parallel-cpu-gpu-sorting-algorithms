// Test comparealgs
package main

import (
	"testing"
)

// TestMaxNumDigits checks maxNumDigits with a multitude of input arrays.
func TestMaxNumDigits(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7}, 1},
		{[]int{7, 6, 5, 45, 3, 26, 1, 10}, 2},
		{[]int{999, 4, 295, 666, 43, 66, 6, 576}, 3},
	}

	for _, c := range cases {
		got := maxNumDigits(c.in)

		if got != c.want {
			t.Errorf("maxNumDigits (%v) == %d, want %d", c.in, got, c.want)
		}
	}
}
