package main

import (
	"fmt"
	go_concurrent_pipeline "github.com/ankit-arora/go-utils/go-concurrent-pipeline"
	go_pipeline "github.com/ankit-arora/go-utils/go-concurrent-pipeline/go-pipeline"
	goroutines_pool "github.com/ankit-arora/go-utils/go-concurrent-pipeline/goroutines-pool"
	go_shutdown_hook "github.com/ankit-arora/go-utils/go-shutdown-hook"
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

	go_shutdown_hook.ADD(func() {
		fmt.Println("stopping pool p")
		p.Shutdown()
		fmt.Println("stopping pool p-done")
	})

	for i := 0; i <= 10; i++ {
		j := i
		p.Submit(func() {
			fmt.Println("task " + strconv.Itoa(j))
		})
	}

	//only pool example ends

	//only pipeline example starts
	pipeline, err := go_pipeline.NewPipeline(10, 2*time.Second, func(key string, data []interface{}) {
		fmt.Print(key + " -> ")
		fmt.Println(data)
	})

	go_shutdown_hook.ADD(func() {
		fmt.Println("shutting down pipeline")
		pipeline.Shutdown()
		fmt.Println("shutting down pipeline-done")
	})

	for i := 0; i < 52; i++ {
		pipeline.Add("0-pipeline", i)
		pipeline.Add("1-pipeline", i*2)
	}

	pipeline.Add("2-pipeline", "foo")
	pipeline.Add("2-pipeline", "bar")

	//time.Sleep(1 * time.Second)

	go func() {
		pipeline.Add("2-pipeline", "foo-1")
	}()
	pipeline.Add("2-pipeline", "bar-1")

	//time.Sleep(10 * time.Second)

	//only pipeline example ends

	//concurrent pipeline example starts
	cp, err := go_concurrent_pipeline.NewConcurrentPipeline(10, 10, 3*time.Second,
		func(key int, i []string) {
			fmt.Printf("%d->%s\n", key, i)
		})

	go_shutdown_hook.ADD(func() {
		fmt.Println("shutting down concurrent_pipeline")
		cp.Shutdown()
		fmt.Println("shutting down concurrent_pipeline-done")
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 51; i++ {
		cp.Add(10, fmt.Sprintf("%d", i))
	}

	time.Sleep(10 * time.Second)

	for i := 51; i < 101; i++ {
		cp.Add(11, fmt.Sprintf("%d", i))
	}

	for i := 101; i < 151; i++ {
		cp.Add(10, fmt.Sprintf("%d", i))
	}

	////concurrent pipeline example ends

	go_shutdown_hook.Wait()
}
