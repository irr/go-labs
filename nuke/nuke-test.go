package main

import (
	"github.com/ortuman/nuke"
)

type Foo struct{ A int }

func main() {
	// Initialize a new monotonic arena with a buffer size of 256KB
	// and a max memory size of 20MB.
	arena := nuke.NewMonotonicArena(256*1024, 80)

	// Allocate a new object of type Foo.
	fooRef := nuke.New[Foo](arena)
	fooRef.A = 42

	// Allocate a Foo slice with a capacity of 10 elements.
	fooSlice := nuke.MakeSlice[Foo](arena, 0, 10)

	// Append 20 elements to the slice allocating
	// the required extra memory from the arena.
	for i := range 20 {
		fooSlice = nuke.SliceAppend(arena, fooSlice, Foo{A: i})
	}

	// ...

	// When done, reset the arena (releasing monotonic buffer memory).
	arena.Reset(true)

	// From here on, any arena reference is invalid.
	// ...
}
