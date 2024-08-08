package main

import (
	"fmt"
	"sync"
)

/*
Data Race:
 Data race occurs when two or more goroutine tries to access a shared variable or object or file.
 run the code as go run -race <file_name>
*/

var sharedVariable string

func updateString(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	sharedVariable = s
}

func main() {

	var wg sync.WaitGroup

	wg.Add(2)

	sharedVariable = "I am a varibale to be shared among go rotuines"

	go updateString("I am first go routine to acess and update the shared variable", &wg)
	go updateString("I am second go routine to acess and update the shared variable", &wg)
	wg.Wait()
	fmt.Println(sharedVariable)

}
