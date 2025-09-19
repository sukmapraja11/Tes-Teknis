package main

import (
	"fmt"
	"math"
)

// Fungsi factorial
func factorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	result := uint64(1)
	for i := 1; i <= n; i++ {
		result *= uint64(i)
	}
	return result
}

// Fungsi f(n) = (n!) / (2^n), hasil dibulatkan ke atas
func f(n int) uint64 {
	num := float64(factorial(n))
	den := math.Pow(2, float64(n))
	result := math.Ceil(num / den)
	return uint64(result)
}

func main() {
	for i := 0; i <= 10; i++ {
		fmt.Printf("f(%d) = %d\n", i, f(i))
	}
}
