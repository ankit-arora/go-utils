package main

import (
	"fmt"
	"goroutines-pool/go-pipeline"
	"goroutines-pool/goroutines-pool"
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
	p.Start()

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

}
