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
	pipeline, err := go_pipeline.NewPipeline(10, 1*time.Second, func(data []interface{}) {
		fmt.Println(data)
	})

	for i := 0; i < 52; i++ {
		pipeline.Add(1)
	}

	//time.Sleep(1 * time.Minute)

	pipeline.Shutdown()
	//only pipeline example ends

	cp, err := go_concurrent_pipeline.NewConcurrentPipeline(3, 20, 10*time.Second, func(i []interface{}) {
		fmt.Println(i)
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 1004; i++ {
		cp.Add(2)
	}
	cp.Shutdown()

}
