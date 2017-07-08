package go_concurrent_pipeline

import (
	"go-utils/go-concurrent-pipeline/go-pipeline"
	"go-utils/go-concurrent-pipeline/goroutines-pool"
	"time"
)

type ConcurrentPipeline struct {
	pipeline *go_pipeline.Pipeline
	pool     *goroutines_pool.Pool
}

func (cp *ConcurrentPipeline) Shutdown() {
	cp.pipeline.Shutdown()
	cp.pool.Stop()
}

func (cp *ConcurrentPipeline) Add(d interface{}) {
	cp.pipeline.Add(d)
}

func NewConcurrentPipeline(maxConcurrency int, pipelineSize int, timeout time.Duration, f func([]interface{})) (*ConcurrentPipeline, error) {
	cp := &ConcurrentPipeline{}
	pool, err := goroutines_pool.NewPool(maxConcurrency)
	if err != nil {
		return nil, err
	}
	cp.pool = pool

	cp.pipeline, err = go_pipeline.NewPipeline(pipelineSize, timeout, func(data []interface{}) {
		temp := make([]interface{}, len(data))
		copy(temp, data)
		cp.pool.Submit(func() {
			f(temp)
		})
	})

	return cp, nil
}
