package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizza = 10

var pizzaMade, pizzaFailed, total int

type Producer struct {
	data chan pizzaOrder
	quit chan chan error
}

type pizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func main() {

	//seed the random number generator
	rand.New(rand.NewSource(time.Now().Unix()))
	color.Cyan("The Pizzeria is open for buisness")
	pizzaProducer := Producer{
		data: make(chan pizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(&pizzaProducer)

	for v := range pizzaProducer.data {
		if v.pizzaNumber <= NumberOfPizza {
			if v.success {
				color.Green(v.message)
				color.Green("order #%d is out for delivery", v.pizzaNumber)
			} else {
				color.Red(v.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done Making Pizza")
			err := pizzaProducer.Close()
			if err != nil {
				color.Red("*** Error closing the channel")
			}
		}
	}

}

func pizzeria(pizzaMaker *Producer) {
	var i = 0
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:
			case quit := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quit)
				return
			}
		}
	}
}

func makePizza(pizzaNumber int) *pizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizza {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzaFailed++
		} else {
			pizzaMade++
		}
		total++

		fmt.Printf("Making pizza #%d. I will take %d  seconds.....\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza Order #%d", pizzaNumber)
		}

		p := pizzaOrder{
			pizzaNumber: pizzaNumber,
			success:     success,
			message:     msg,
		}

		return &p

	}

	return &pizzaOrder{pizzaNumber: pizzaNumber}
}
