package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	const numberOfSubroutines = 4
	var x int
	channels := make([]chan int, numberOfSubroutines)
	var integerArray []int
	for {
		fmt.Println("Type the next one . If you want to start sorting, type anything else than integer")
		_, err := fmt.Scanf("%d", &x)
		if err != nil {
			break
		}
		integerArray = append(integerArray, x)
	}
	sortedArray := make([]int, len(integerArray))

	if len(integerArray) == 0 {
		fmt.Println("You need to type at least one integer")
	}

	wg := sync.WaitGroup{}
	for i := 0; i < numberOfSubroutines; i++ {
		wg.Add(1)
		if i == numberOfSubroutines-1 {
			go sortArray(integerArray[i*(len(integerArray)/numberOfSubroutines):], &wg, channels[i])
		} else {
			go sortArray(integerArray[i*(len(integerArray)/numberOfSubroutines):(i+1)*(len(integerArray)/numberOfSubroutines)], &wg, channels[i])
		}
	}

	wg.Wait()
	for i := 0; i < numberOfSubroutines; i++ {
		select {
		case t, ok := <-channels[i]:
			if ok {
				sortedArray = append(sortedArray, t)
			} else {
				break
			}
		default:
			fmt.Println("No value ready, moving on.")
		}
	}
	fmt.Println(sortedArray)
}

func sortArray(arrayOfIntegers []int, wg *sync.WaitGroup, channel chan int) {
	fmt.Println("I will sort:", arrayOfIntegers)
	defer func() {
		close(channel)
		wg.Done()
	}()
	sort.Ints(arrayOfIntegers)
	for _, i := range arrayOfIntegers {
		channel <- i
	}
	fmt.Println("I sorted:", arrayOfIntegers)
}
