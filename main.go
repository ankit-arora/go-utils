package main

import (
	"fmt"
	"goroutines-pool/goroutines-pool"
	"strconv"
)

func main() {
	p, err := goroutines_pool.NewPool(5)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.Start()

	for i := 0; i <= 10; i++ {
		j := i
		p.Submit(func() {
			fmt.Println("task " + strconv.Itoa(j))
		})
	}

	p.Stop()
}
