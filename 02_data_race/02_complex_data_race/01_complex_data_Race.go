package main

import (
	"fmt"
	"sync"
)

type Income struct {
	Source string
	Amount float64
}

func main() {

	var bankBalance float64
	var wg sync.WaitGroup
	var mutex sync.Mutex
	fmt.Println("Initial bank balance before computation is:", bankBalance)

	incomes := []Income{
		{Source: "Main Job", Amount: 3.01},
		{Source: "Teaching", Amount: 4.01},
		{Source: "Part-time", Amount: 3.225},
		{Source: "Real-estate", Amount: 5.255},
		{Source: "Freelancing", Amount: 5.155},
	}

	wg.Add(len(incomes))

	for i, v := range incomes {
		go CalculateYearIncome(i, v, &bankBalance, &wg, &mutex)
	}
	wg.Wait()
	fmt.Println("Final Bank Balance for a year is INR:", bankBalance)

}

func CalculateYearIncome(i int, income Income, bankBalance *float64, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	for week := 1; week <= 52; week++ {
		mutex.Lock()
		temp := *bankBalance
		temp += income.Amount
		*bankBalance = temp
		mutex.Unlock()

		fmt.Printf("On week %d, you earned INR %v from %s\n", week, income.Amount, income.Source)
	}

}
