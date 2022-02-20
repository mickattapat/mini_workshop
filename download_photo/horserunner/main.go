package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func running(name string, track chan struct{}) {
	defer func() {
		fmt.Printf("%s stopped running\n", name)
	}()
	rand.Seed(time.Now().UnixNano())
	chFinish := time.After(time.Duration(1*time.Second + time.Duration(rand.Intn(100))*time.Millisecond))
	track <- struct{}{}

	select {
	case <-chFinish:
		track <- struct{}{}
	}
}

func main() {
	defer func() {
		fmt.Println("Done")
	}()

	track1 := make(chan struct{})
	track2 := make(chan struct{})
	track3 := make(chan struct{})
	abort := make(chan struct{})

	go running("h1", track1)
	go running("h2", track2)
	go running("h3", track3)
	lapsChan := make(chan time.Time)
	go func() {
		labsTricker := time.NewTicker(1 * time.Second)
		for {
			lapsChan <- (<-labsTricker.C)
		}

	}()
	go func() {
		os.Stdin.Read(make([]byte, 1))
		fmt.Println("About for abort")
		abort <- struct{}{}
	}()
	for done := false; !done; {
		select {
		case <-track1:
			fmt.Println("Winner Horse 1 ")
			done = true
		case <-track2:
			fmt.Println("Winner Horse 2 ")
			done = true
		case <-track3:
			fmt.Println("Winner Horse 3 ")
			done = true
		case v := <-lapsChan:
			fmt.Println("Horses are running : ", v.Second())
		case <-abort:
			done = true
			break
		}

	}
	fmt.Println("Done !!")
}
