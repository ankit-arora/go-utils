package main

import (
	"fmt"
	"go-utils/go-concurrent-pipeline/go-pipeline"
	"go-utils/go-concurrent-pipeline/goroutines-pool"

	"go-utils/go-concurrent-pipeline"
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
	pipeline, err := go_pipeline.NewPipeline(10, 2*time.Second, func(key string, data []interface{}) {
		fmt.Print(key + " -> ")
		fmt.Println(data)
	})

	for i := 0; i < 52; i++ {
		pipeline.Add("0-pipeline", i)
		pipeline.Add("1-pipeline", i*2)
	}

	pipeline.Add("2-pipeline", "foo")
	pipeline.Add("2-pipeline", "bar")

	time.Sleep(4 * time.Second)

	pipeline.Add("2-pipeline", "foo-1")
	pipeline.Add("2-pipeline", "bar-1")

	time.Sleep(4 * time.Second)

	pipeline.Shutdown()
	//only pipeline example ends

	//concurrent pipeline example starts
	cp, err := go_concurrent_pipeline.NewConcurrentPipeline(10, 10, 3*time.Second,
		func(key string, i []interface{}) {
			fmt.Print(key + " -> ")
			fmt.Println(i)
		})

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 51; i++ {
		cp.Add("0", i)
	}

	//time.Sleep(10 * time.Second)

	for i := 51; i < 101; i++ {
		cp.Add("1", i)
	}

	for i := 101; i < 151; i++ {
		cp.Add("0", i)
	}

	cp.Shutdown()
	//concurrent pipeline example ends
}
