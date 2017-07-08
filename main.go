package main

import (
	"fmt"
	"go-utils/go-concurrent-pipeline"
	"go-utils/go-concurrent-pipeline/go-pipeline"
	"go-utils/go-concurrent-pipeline/goroutines-pool"

	"strconv"
	"time"
)

func main() {
	//only pool example starts
	p, err := goroutines_pool.NewPool(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i <= 10; i++ {
		j := i
		p.Submit(func() {
			fmt.Println("task " + strconv.Itoa(j))
		})
	}

	p.Stop()
	//only pool example ends

	//only pipeline example starts
	pipeline, err := go_pipeline.NewPipeline(10, 10*time.Second, func(data []interface{}) {
		fmt.Println(data)
	})

	for i := 0; i < 52; i++ {
		pipeline.Add(1)
	}

	//time.Sleep(1 * time.Minute)

	pipeline.Shutdown()
	//only pipeline example ends

	//concurrent pipeline example starts
	cp, err := go_concurrent_pipeline.NewConcurrentPipeline(10, 21, 10*time.Second,
		func(i []interface{}) {
			fmt.Println(i)
		})

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 510; i++ {
		cp.Add(2)
	}
	cp.Shutdown()
	//concurrent pipeline example ends
}
