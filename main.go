package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// Avoid concurrent map writes with a mutex
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}
	data := map[string]int{"zero": 0, "one": 1}
	z := 0
	wg.Add(3)
	go func(m *sync.Mutex, wg *sync.WaitGroup) {
		start := time.Now()
		for i := 0; i < 1000000; i++ {
			m.Lock()
			z = data["zero"]
			m.Unlock()
		}
		fmt.Println(time.Since(start))
		wg.Done()
	}(m, wg)
	go func(m *sync.Mutex, wg *sync.WaitGroup) {
		start := time.Now()
		for i := 0; i < 1000000; i++ {
			m.Lock()
			z = data["one"]
			m.Unlock()
		}
		fmt.Println(time.Since(start))
		wg.Done()
	}(m, wg)
	go func(m *sync.Mutex, wg *sync.WaitGroup) {
		for i := 0; i < 1000000; i++ {
			m.Lock()
			data["zero"] = 2
			m.Unlock()
		}
		wg.Done()
	}(m, wg)
	wg.Wait()
	fmt.Println(z)

	// Allow concurrent map reads with a read-write mutex (actually slower unless there are many more times as reads than writes)
	wg = &sync.WaitGroup{}
	rwm := &sync.RWMutex{}
	data = map[string]int{"zero": 0, "one": 1}
	wg.Add(3)
	go func(m *sync.RWMutex, wg *sync.WaitGroup) {
		start := time.Now()
		for i := 0; i < 1000000; i++ {
			m.RLock()
			z = data["zero"]
			m.RUnlock()
		}
		fmt.Println(time.Since(start))
		wg.Done()
	}(rwm, wg)
	go func(m *sync.RWMutex, wg *sync.WaitGroup) {
		start := time.Now()
		for i := 0; i < 1000000; i++ {
			m.RLock()
			z = data["one"]
			m.RUnlock()
		}
		fmt.Println(time.Since(start))
		wg.Done()
	}(rwm, wg)
	go func(m *sync.RWMutex, wg *sync.WaitGroup) {
		for i := 0; i < 1000000; i++ {
			m.Lock()
			data["zero"] = 2
			m.Unlock()
		}
		wg.Done()
	}(rwm, wg)
	wg.Wait()
	fmt.Println(z)

	// Get a buffered default value from a channel that was closed
	wg = &sync.WaitGroup{}
	ch := make(chan int, 1)
	wg.Add(2)
	go func(ch <-chan int, wg *sync.WaitGroup) {
		if msg, ok := <-ch; ok {
			fmt.Println(msg)
		}
		if msg, ok := <-ch; ok {
			fmt.Println(msg)
		}
		wg.Done()
	}(ch, wg)
	go func(ch chan<- int, wg *sync.WaitGroup) {
		ch <- 1
		close(ch)
		wg.Done()
	}(ch, wg)
	wg.Wait()

	// Read all messages from a channel that was closed
	wg = &sync.WaitGroup{}
	ch = make(chan int)
	wg.Add(2)
	go func(ch <-chan int, wg *sync.WaitGroup) {
		for msg := range ch {
			fmt.Println(msg)
		}
		wg.Done()
	}(ch, wg)
	go func(ch chan<- int, wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}(ch, wg)
	wg.Wait()

	// Select whichever channel produces a value first
	wg = &sync.WaitGroup{}
	ch = make(chan int)
	ch2 := make(chan int)
	wg.Add(3)
	go func(ch <-chan int, wg *sync.WaitGroup) {
		select {
		case msg := <-ch:
			fmt.Println(msg)
			<-ch2
		case msg := <-ch2:
			fmt.Println(msg)
			<-ch
		}
		wg.Done()
	}(ch, wg)
	go func(ch chan<- int, wg *sync.WaitGroup) {
		time.Sleep(time.Duration(rand.Intn(100)))
		ch <- 1
		close(ch)
		wg.Done()
	}(ch, wg)
	go func(ch chan<- int, wg *sync.WaitGroup) {
		time.Sleep(time.Duration(rand.Intn(100)))
		ch2 <- 2
		close(ch2)
		wg.Done()
	}(ch, wg)
	wg.Wait()

	// Select a default value when no channel produces a message
	ch = make(chan int)
	ch2 = make(chan int)
	select {
	case msg := <-ch:
		fmt.Println(msg)
		<-ch2
	case msg := <-ch2:
		fmt.Println(msg)
		<-ch
	default:
		fmt.Println("default")
	}
}
