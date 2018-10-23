package main

import (
	"fmt"
	"golang.org/x/tour/tree"
	"time"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	ch <- t.Value
	if tl := t.Left; tl != nil {
		Walk(tl, ch)
	}
	if tr := t.Left; tr != nil {
		Walk(tr, ch)
	}
}

func fibonacci(n int, ch chan int) {
	defer close(ch)
	x, y := 0, 1
	for i := 0; i < n; i++ {
		ch <- x
		time.Sleep(time.Second)
		x, y = y, x+y
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
// func Same(t1, t2 *tree.Tree) bool
func main() {
	t := tree.New(1)
	fmt.Println(t)
	ch := make(chan int)
	Walk(t, ch)
	//	go fibonacci(10, ch)
	go func() {
		for i := range ch {
			fmt.Printf("% d", i)
		}
	}()
	fmt.Printf("\n")
}
