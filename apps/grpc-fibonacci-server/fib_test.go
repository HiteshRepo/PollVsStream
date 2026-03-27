package grpcfibonacciserver

import "testing"

func TestFib(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{5, 5},
		{10, 55},
		{15, 610},
	}

	for _, tt := range tests {
		got := fib(tt.n)
		if got != tt.want {
			t.Errorf("fib(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
