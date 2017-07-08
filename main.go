package main

import (
	"fmt"
	"goroutines-pool/goroutines-pool"
)

func main() {
	p, err := goroutines_pool.NewPool(3)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.Start()
	p.Submit(func() {
		fmt.Println("task 1")
	})
	p.Submit(func() {
		fmt.Println("task 2")
	})
	p.Submit(func() {
		fmt.Println("task 3")
	})
	p.Stop()
}
