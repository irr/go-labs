// go test -v

package main

import "testing"

func Sqrt(x float64) float64 {
    z := 0.0
    for i := 0; i < 1000; i++ {
        z -= (z*z - x) / (2 * x)
    }
    return z
}

func TestSqrt(t *testing.T) {
    const in, out = 4, 2
    if x := Sqrt(in); x != out {
        t.Errorf("Sqrt(%v) = %v, want %v", in, x, out)
    }
}
