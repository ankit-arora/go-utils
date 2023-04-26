package go_concurrent_pipeline

import (
	go_pipeline "github.com/ankit-arora/go-utils/go-concurrent-pipeline/go-pipeline"
	goroutines_pool "github.com/ankit-arora/go-utils/go-concurrent-pipeline/goroutines-pool"
	"time"
)

type ConcurrentPipeline[K comparable, T any] struct {
	pipeline *go_pipeline.Pipeline[K, T]
	pool     *goroutines_pool.Pool
}

func (cp *ConcurrentPipeline[K, T]) Shutdown() {
	cp.pipeline.Shutdown()
	cp.pool.Shutdown()
}

func (cp *ConcurrentPipeline[K, T]) Add(key K, i T) {
	cp.pipeline.Add(key, i)
}

func NewConcurrentPipeline[K comparable, T any](maxConcurrency int, pipelineSize int, timeout time.Duration, f func(K, []T)) (*ConcurrentPipeline[K, T], error) {
	cp := &ConcurrentPipeline[K, T]{}
	pool, err := goroutines_pool.NewPool(maxConcurrency)
	if err != nil {
		return nil, err
	}
	cp.pool = pool

	cp.pipeline, err = go_pipeline.NewPipeline(pipelineSize, timeout, func(key K, i []T) {
		temp := make([]T, len(i))
		copy(temp, i)
		cp.pool.Submit(func() {
			f(key, temp)
		})
	})

	return cp, nil
}
