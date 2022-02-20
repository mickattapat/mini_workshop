package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const N = 1000000
	add := 0 // n
	sub := 0 // -n
	mu := sync.Mutex{}
	for i := 0; i < N; i++ {
		add++
		go func() {
			defer func() {
				mu.Lock()   // Lock
				sub--       // critical section
				mu.Unlock() //	Unlock
			}()
		}()
	}
	for {
		h, m, s := time.Now().Clock()
		fmt.Printf("%d-%d-%d : add %d , sub %d\r", h, m, s, add, sub)
		if add == N && sub == -N {
			fmt.Println("success")
			break
		}
	}
	fmt.Println("main done", add, sub)
}
