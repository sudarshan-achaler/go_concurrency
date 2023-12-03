package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// MaxPizzas represents the maximum number of pizzas that can be in the order queue.
var MaxPizzas int = 10

// wg is used to wait for goroutines to finish.
var wg sync.WaitGroup

// orders is a buffered channel to represent the pizza order queue.
var orders chan string = make(chan string, MaxPizzas)

// stopSignal is a channel used to signal termination to the chefs and consumers.
var stopSignal chan int = make(chan int)

// pizzaMenu represents the types of pizzas available.
var pizzaMenu []string = []string{"Margherita", "Pepperoni", "Vegetarian", "Hawaiian"}

// cookPizza simulates a chef cooking pizzas.
func cookPizza(chefID int) {
	fmt.Printf("chef %v started making pizzas\n", chefID)
	defer wg.Done()

	for {
		select {
		case <-stopSignal:
			fmt.Printf("received stop signal, chef %v exiting...\n", chefID)
			return
		default:
			// Simulate time taken to cook a pizza
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			order := pizzaMenu[rand.Intn(len(pizzaMenu))]
			orders <- order // Send the cooked pizza to the order queue
			fmt.Printf("%v pizza cooked and ready to consume\n", order)
		}
	}
}

// orderPizza simulates a consumer ordering and consuming pizzas.
func orderPizza(consumerID int) {
	fmt.Printf("consumer %v started consuming pizzas\n", consumerID)
	defer wg.Done()

	for {
		// Simulate time taken to consume a pizza
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

	// Start pizza producers (chefs)
	for i := 0; i < pizzaProducers; i++ {
		wg.Add(1)
		go cookPizza(i)
	}

	// Start pizza consumers
	for j := 0; j < pizzaConsumers; j++ {
		wg.Add(1)
		go orderPizza(j)
	}

	// Simulate a period of time for chefs and consumers to work
	time.Sleep(20 * time.Second)

	// Send stop signal to chefs and consumers
	close(stopSignal)

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the orders channel
	close(orders)
}
