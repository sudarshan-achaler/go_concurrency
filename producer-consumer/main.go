package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var MaxPizzas int = 10
var wg sync.WaitGroup
var orders chan string = make(chan string, MaxPizzas)
var stopSignal chan int = make(chan int)

var pizzaMenu []string = []string{"Margherita", "Pepperoni", "Vegetarian", "Hawaiian"}

func cookPizza(chefID int) {
	fmt.Printf("chef %v started making pizzas\n", chefID)
	defer wg.Done()

	for {
		select {
		case <-stopSignal:
			fmt.Printf("received stop signal, chef %v exiting...\n", chefID)
			return
		default:
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			order := pizzaMenu[rand.Intn(len(pizzaMenu))]
			orders <- order
			fmt.Printf("%v pizza cooked and ready to consume\n", order)
		}
	}

}

func orderPizza(consumerID int) {
	fmt.Printf("consumer %v started consuming pizzas\n", consumerID)
	defer wg.Done()

	for {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		select {
		case <-stopSignal:
			fmt.Printf("received stop signal, consumer %v exiting...\n", consumerID)
			return
		case order := <-orders:
			fmt.Printf("%v pizza consumed\n", order)
		}
	}

}

func main() {
	pizzaProducers := 2
	pizzaConsumers := 3

	for i := 0; i < pizzaProducers; i++ {
		wg.Add(1)
		go cookPizza(i)
	}

	for j := 0; j < pizzaConsumers; j++ {
		wg.Add(1)
		go orderPizza(j)
	}

	time.Sleep(20 * time.Second)

	close(stopSignal)
	wg.Wait()
	close(orders)

}
