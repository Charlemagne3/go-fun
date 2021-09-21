package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
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
