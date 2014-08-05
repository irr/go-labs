package main

import (
    "math/rand"
    "fmt"
    "time"
)

type Vector []float64 

func (v Vector) MeanAndSignum() (mean float64, signum []int) { 
    total := 0.0
    signum = make([]int, len(v)) 
    for i, value := range v { 
        total += value
        switch { 
            case value < 0.0: signum[i] = -1 
            case value == 0.0: signum[i] = 0 
            case value > 0.0: signum[i] = 1 } 
        } 
    mean = total / float64(len(v)) 
    return
}

func main() {
    rand.Seed( time.Now().UTC().UnixNano())
    
    v := make(Vector, 10)
    for i := 0; i < len(v); i++ {
        v[i] = rand.Float64()
        if (v[i] > 0.5) {
            v[i] *= -1
        }
    }
    mean, signum := v.MeanAndSignum()
    fmt.Printf("mean = %v\n", mean)

    for i, e := range v {
        fmt.Printf("%4d = %v [%v]\n", i, e, signum[i])
    }
}