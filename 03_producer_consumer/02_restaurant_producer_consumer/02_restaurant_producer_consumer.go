package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizza = 20
const NumberOfBurger = 20

var pizzaMade, pizzaFailed, total int

type Producer struct {
	pizzadata  chan pizzaOrder
	burgerdata chan burgerOrder
	quit       chan chan error
}

type pizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

type burgerOrder struct {
	burgerNumber int
	message      string
	success      bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func main() {
	var wg sync.WaitGroup
	//seed the random number generator
	rand.New(rand.NewSource(time.Now().Unix()))
	color.Cyan("The Restaurant is open for buisness")
	Produce := Producer{
		pizzadata:  make(chan pizzaOrder),
		burgerdata: make(chan burgerOrder),
		quit:       make(chan chan error),
	}
	wg.Add(4)

	go pizzeria(&Produce, &wg)
	go burger(&Produce, &wg)

	go func(producer *Producer) {
		defer wg.Done()
		for v := range producer.pizzadata {
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
				err := producer.Close()
				if err != nil {
					color.Red("*** Error closing the channel")
				}
			}
		}
	}(&Produce)

	go func(producer *Producer) {
		defer wg.Done()

		for v := range producer.burgerdata {
			if v.burgerNumber <= NumberOfBurger {
				if v.success {
					color.Green(v.message)
					color.Green("order #%d is out for delivery", v.burgerNumber)
				} else {
					color.Red(v.message)
					color.Red("The customer is really mad!")
				}
			} else {
				color.Cyan("Done Making Pizza")
				err := producer.Close()
				if err != nil {
					color.Red("*** Error closing the channel")
				}
			}
		}
	}(&Produce)

	wg.Wait()

}

func burger(burgermaker *Producer, wg *sync.WaitGroup) {
	defer wg.Done()
	var i = 0
	for {
		currentBurger := makeBurger(i)
		if currentBurger != nil {
			i = currentBurger.burgerNumber

			select {
			case burgermaker.burgerdata <- *currentBurger:

			case quit := <-burgermaker.quit:
				close(burgermaker.burgerdata)
				close(quit)
				return
			}
		}
	}

}

func pizzeria(pizzaMaker *Producer, wg *sync.WaitGroup) {
	defer wg.Done()
	var i = 0
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.pizzadata <- *currentPizza:
			case quit := <-pizzaMaker.quit:
				close(pizzaMaker.pizzadata)
				close(quit)
				return
			}
		}
	}
}

func makeBurger(burgerNumber int) *burgerOrder {

	burgerNumber++
	if burgerNumber <= NumberOfBurger {
		delay := rand.Intn(5) + 1
		fmt.Println("Received order #%d!", burgerNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzaFailed++
		} else {
			pizzaMade++
		}
		total++

		fmt.Printf("Making Burger #%d. I will take %d  seconds.....\n", burgerNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for Burger #%d!", burgerNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making Burger #%d!", burgerNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Burger Order #%d", burgerNumber)
		}

		p := burgerOrder{
			burgerNumber: burgerNumber,
			success:      success,
			message:      msg,
		}

		return &p
	}

	return &burgerOrder{
		burgerNumber: burgerNumber,
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
