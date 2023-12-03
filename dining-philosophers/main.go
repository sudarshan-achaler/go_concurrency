package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	Name      string
	LeftFork  int
	RightFork int
}

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

	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for _, philosopher := range philosophers {
		go dine(philosopher, wg, seated, forks)
	}

	wg.Wait()

}

func dine(philosopher Philosopher, wg *sync.WaitGroup, seated *sync.WaitGroup, forks map[int]*sync.Mutex) {
	defer wg.Done()

	fmt.Printf("%v is seated on table\n", philosopher.Name)
	seated.Done()

	seated.Wait()

	for i := hunger; i > 0; i-- {
		if philosopher.LeftFork < philosopher.RightFork {
			forks[philosopher.LeftFork].Lock()
			fmt.Printf("%v has picked left fork %v\n", philosopher.Name, philosopher.LeftFork)
			forks[philosopher.RightFork].Lock()
			fmt.Printf("%v has picked right fork %v\n", philosopher.Name, philosopher.RightFork)
		} else {
			forks[philosopher.RightFork].Lock()
			fmt.Printf("%v has picked right fork %v\n", philosopher.Name, philosopher.RightFork)
			forks[philosopher.LeftFork].Lock()
			fmt.Printf("%v has picked left fork %v\n", philosopher.Name, philosopher.LeftFork)
		}

		time.Sleep(eatTime)

		forks[philosopher.RightFork].Unlock()
		forks[philosopher.LeftFork].Unlock()

		fmt.Printf("%v has put down his forks\n", philosopher.Name)

	}

	fmt.Printf("%v is satisfied and left the table\n", philosopher.Name)
}
