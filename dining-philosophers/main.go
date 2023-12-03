package main

import (
	"fmt"
	"sync"
	"time"
)

// Philosopher represents a philosopher with a name and the indices of their left and right forks.
type Philosopher struct {
	Name      string
	LeftFork  int
	RightFork int
}

// Global variables
var philosophers = []Philosopher{
	{Name: "A", LeftFork: 4, RightFork: 0},
	{Name: "B", LeftFork: 0, RightFork: 1},
	{Name: "C", LeftFork: 1, RightFork: 2},
	{Name: "D", LeftFork: 2, RightFork: 3},
	{Name: "E", LeftFork: 3, RightFork: 4},
}

var eatTime = 2 * time.Second
var hunger = 3

func main() {
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("Table is empty")

	// WaitGroup to wait for all philosophers to finish dining
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// WaitGroup to ensure all philosophers are seated before they start dining
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// Map to represent forks, each fork is protected by a mutex
	forks := make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// Start a goroutine for each philosopher
	for _, philosopher := range philosophers {
		go dine(philosopher, wg, seated, forks)
	}

	// Wait for all philosophers to finish dining
	wg.Wait()
}

// dine represents the dining behavior of a philosopher
func dine(philosopher Philosopher, wg *sync.WaitGroup, seated *sync.WaitGroup, forks map[int]*sync.Mutex) {
	defer wg.Done()

	// Philosopher is seated
	fmt.Printf("%v is seated on the table\n", philosopher.Name)
	seated.Done()

	// Wait for all philosophers to be seated
	seated.Wait()

	// Dining loop
	for i := hunger; i > 0; i-- {
		// Determine the order to pick up forks based on indices
		if philosopher.LeftFork < philosopher.RightFork {
			// Pick up left fork first
			forks[philosopher.LeftFork].Lock()
			fmt.Printf("%v has picked left fork %v\n", philosopher.Name, philosopher.LeftFork)
			forks[philosopher.RightFork].Lock()
			fmt.Printf("%v has picked right fork %v\n", philosopher.Name, philosopher.RightFork)
		} else {
			// Pick up right fork first
			forks[philosopher.RightFork].Lock()
			fmt.Printf("%v has picked right fork %v\n", philosopher.Name, philosopher.RightFork)
			forks[philosopher.LeftFork].Lock()
			fmt.Printf("%v has picked left fork %v\n", philosopher.Name, philosopher.LeftFork)
		}

		// Simulate eating
		time.Sleep(eatTime)

		// Release forks
		forks[philosopher.RightFork].Unlock()
		forks[philosopher.LeftFork].Unlock()

		// Philosopher puts down the forks
		fmt.Printf("%v has put down his forks\n", philosopher.Name)
	}

	// Philosopher is satisfied and leaves the table
	fmt.Printf("%v is satisfied and left the table\n", philosopher.Name)
}
