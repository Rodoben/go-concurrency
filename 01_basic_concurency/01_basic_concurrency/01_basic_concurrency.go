package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"epsilon",
		"zeta",
		"thetha",
		"ronald",
	}

	wg.Add(len(words))
	for _, v := range words {
		go printSomething(v, &wg)
	}
	wg.Wait()
	wg.Add(1)

	printSomething("Second statement to be printed!", &wg)
	wg.Wait()
}

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(s)
}
