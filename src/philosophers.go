package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const numOfPhilosophers = 5
const numOfChopsticks = 5
const maxTimesEating = 3
const maxNumberOfPhilosophersEatingConcurrently = 2

var philosophersEatingConcurrently = 0

type philosopher struct {
	rChopstick *chopstick
	lChopstick *chopstick
	timesAte   int
}

func (p *philosopher) eat(number int) {
	p.takeChopsticks()
	fmt.Println("starting to eat", number)
	p.timesAte++
	fmt.Println("finishing eating", number)
	p.returnChopsticks()
}

func (p *philosopher) takeChopsticks() {
	//each philosopher takes randomly right or left chopstick
	randomChopstick := rand.Intn(2)
	if randomChopstick == 0 {
		p.lChopstick.Lock()
		p.rChopstick.Lock()
	} else {
		p.rChopstick.Lock()
		p.lChopstick.Lock()
	}
}
func (p *philosopher) returnChopsticks() {
	p.lChopstick.Unlock()
	p.rChopstick.Unlock()
}

func (p *philosopher) dine(wg *sync.WaitGroup, m *sync.Mutex, requestPermissionChannel chan<- string, grantPermissionChannel <-chan int) {
	//each philosopher dines. The philosopher requests permission by sending message through the requestPermissionChannel to the host
	//If the permission is granted by the host , he then eats (using the locks of the chopsticks). If not , he requests again
	for p.timesAte < maxTimesEating {
		requestPermissionChannel <- "I'm hungry"
		philosopherNumber := <-grantPermissionChannel
		if philosopherNumber == -1 {
			continue
		}
		m.Lock()
		philosophersEatingConcurrently = philosophersEatingConcurrently + 1
		m.Unlock()
		p.eat(philosopherNumber)
		m.Lock()
		philosophersEatingConcurrently = philosophersEatingConcurrently - 1
		m.Unlock()
	}
	wg.Done()
}

type host struct {
}

type chopstick struct {
	sync.Mutex
}

func (h *host) listenAndGrantPermissions(requestPermissionChannels []chan string, grantPermissionChannels []chan int) {
	//the hosts listens all the philosophers . If someone is hungry , the host checks how many are eating concurrently
	//and if they ar less than two he/she responds with the number of the philosopher else he/she declines the permission
	for {
		select {
		case <-requestPermissionChannels[0]:
			if philosophersEatingConcurrently < maxNumberOfPhilosophersEatingConcurrently {
				grantPermissionChannels[0] <- 1
			} else {
				grantPermissionChannels[0] <- -1
			}
		case <-requestPermissionChannels[1]:
			if philosophersEatingConcurrently < maxNumberOfPhilosophersEatingConcurrently {
				grantPermissionChannels[1] <- 2
			} else {
				grantPermissionChannels[1] <- -1
			}
		case <-requestPermissionChannels[2]:
			if philosophersEatingConcurrently < maxNumberOfPhilosophersEatingConcurrently {
				grantPermissionChannels[2] <- 3
			} else {
				grantPermissionChannels[2] <- -1
			}
		case <-requestPermissionChannels[3]:
			if philosophersEatingConcurrently < maxNumberOfPhilosophersEatingConcurrently {
				grantPermissionChannels[3] <- 4
			} else {
				grantPermissionChannels[3] <- -1
			}
		case <-requestPermissionChannels[4]:
			if philosophersEatingConcurrently < maxNumberOfPhilosophersEatingConcurrently {
				grantPermissionChannels[4] <- 5
			} else {
				grantPermissionChannels[4] <- -1
			}
		}
	}
}

func main() {
	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	requestPermissionChannels := make([]chan string, numOfPhilosophers)
	grantPermissionChannels := make([]chan int, numOfPhilosophers)
	philosophers := make([]*philosopher, numOfPhilosophers)
	chopsticks := make([]*chopstick, numOfChopsticks)
	h := host{}

	for i := 0; i < numOfPhilosophers; i++ {
		chopsticks[i] = &chopstick{}
	}
	for i := 0; i < numOfPhilosophers; i++ {
		requestPermissionChannels[i] = make(chan string)
		grantPermissionChannels[i] = make(chan int)

		philosophers[i] = &philosopher{
			lChopstick: chopsticks[i],
			rChopstick: chopsticks[(i+1)%numOfChopsticks],
		}
	}
	// go routine for the host to listen at the request channels and respond to grant permission channels
	go h.listenAndGrantPermissions(requestPermissionChannels, grantPermissionChannels)

	for i := 0; i < numOfPhilosophers; i++ {
		wg.Add(1)
		go philosophers[i].dine(&wg, &m, requestPermissionChannels[i], grantPermissionChannels[i])
	}
	wg.Wait()
}
