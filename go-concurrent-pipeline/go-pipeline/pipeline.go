package go_pipeline

import "time"

type data[T any] struct {
	key string
	i   T
}

type Pipeline[T any] struct {
	pipe         chan data[T]
	pipeSize     int
	timeout      time.Duration
	pipeLineTask func(string, []T)
	shutdown     chan struct{}
	wait         chan struct{}
}

func (p *Pipeline[T]) Add(key string, i T) {
	p.pipe <- data[T]{key, i}
}

func (p *Pipeline[T]) clear(stateMap map[string][]T) {
	for key, state := range stateMap {
		if len(state) > 0 {
			p.pipeLineTask(key, state)
		}
	}
}

func (p *Pipeline[T]) start() {
	go func() {
		var t *time.Timer
		stateMap := make(map[string][]T)
		shutdown := false
		for !shutdown {
			t = time.NewTimer(p.timeout)
			select {
			case d := <-p.pipe:
				key := d.key
				i := d.i
				stateMap[key] = append(stateMap[key], i)
				if len(stateMap[key]) == p.pipeSize {
					//fmt.Println("from data")
					p.pipeLineTask(key, stateMap[key])
					stateMap[key] = nil
				}
			case <-t.C:
				p.clear(stateMap)
				stateMap = make(map[string][]T)
			case <-p.shutdown:
				p.clear(stateMap)
				stateMap = make(map[string][]T)
				shutdown = true
			}
			t.Stop()
		}
		p.wait <- struct{}{}
	}()
}

func (p *Pipeline[T]) Shutdown() {
	close(p.shutdown)
	<-p.wait
	close(p.pipe)
}

func NewPipeline[T any](pipeSize int, timeout time.Duration, f func(string, []T)) (*Pipeline[T], error) {
	p := &Pipeline[T]{}
	p.pipeSize = pipeSize
	p.timeout = timeout
	p.pipeLineTask = f
	p.pipe = make(chan data[T])
	p.shutdown = make(chan struct{})
	p.wait = make(chan struct{})
	p.start()
	return p, nil
}
