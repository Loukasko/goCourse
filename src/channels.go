package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

func main() {
	const numberOfSubroutines = 4

	var integerArray []int
	var sortedArray []int
	for {
		var x int
		fmt.Println("Type a number to sort . If you want to stop type anything other than integer")
		_, err := fmt.Scanf("%d", &x)
		if err != nil {
			break
		}
		integerArray = append(integerArray, x)
	}
	if len(integerArray) == 0 {
		fmt.Println("You need to type at least one integer")
		os.Exit(1)
	}

	var channels [numberOfSubroutines]chan []int
	for i := range channels {
		channels[i] = make(chan []int)
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

	for i := 0; i < numberOfSubroutines; i++ {
		select {
		case t, ok := <-channels[i]:
			if ok {
				sortedArray = append(sortedArray, t...)
			} else {
				break
			}
		}
	}

	wg.Wait()
	sort.Ints(sortedArray)
	fmt.Println(sortedArray)
}

func sortArray(arrayOfIntegers []int, wg *sync.WaitGroup, channel chan []int) {
	fmt.Println("I will sort:", arrayOfIntegers)
	defer func() {
		close(channel)
		wg.Done()
	}()
	sort.Ints(arrayOfIntegers)
	channel <- arrayOfIntegers
}
