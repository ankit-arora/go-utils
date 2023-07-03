package go_pipeline

import "time"

type data[K comparable, T any] struct {
	key K
	i   T
}

type Pipeline[K comparable, T any] struct {
	pipe         chan data[K, T]
	pipeSize     int
	timeout      time.Duration
	pipeLineTask func(K, []T)
	shutdown     chan struct{}
	wait         chan struct{}
}

func (p *Pipeline[K, T]) Add(key K, i T) bool {
	select {
	case <-p.shutdown:
		return false
	default:
		p.pipe <- data[K, T]{key, i}
	}
	return true
}

func (p *Pipeline[K, T]) clear(stateMap map[K][]T) {
	for key, state := range stateMap {
		if len(state) > 0 {
			p.pipeLineTask(key, state)
		}
	}
}

func (p *Pipeline[K, T]) start() {
	go func() {
		var t *time.Timer
		stateMap := make(map[K][]T)
		shutdown := false
		t = time.NewTimer(p.timeout)
		for !shutdown {
			select {
			case <-p.shutdown:
				p.clear(stateMap)
				stateMap = make(map[K][]T)
				shutdown = true
				t.Stop()
			case <-t.C:
				p.clear(stateMap)
				stateMap = make(map[K][]T)
				t.Stop()
				t.Reset(p.timeout)
			case d := <-p.pipe:
				key := d.key
				i := d.i
				stateMap[key] = append(stateMap[key], i)
				if len(stateMap[key]) == p.pipeSize {
					//fmt.Println("from data")
					p.pipeLineTask(key, stateMap[key])
					stateMap[key] = nil
				}
			}
		}
		p.wait <- struct{}{}
	}()
}

func (p *Pipeline[K, T]) Shutdown() {
	close(p.shutdown)
	<-p.wait
	close(p.pipe)
}

func NewPipeline[K comparable, T any](pipeSize int, timeout time.Duration, f func(K, []T)) (*Pipeline[K, T], error) {
	p := &Pipeline[K, T]{}
	p.pipeSize = pipeSize
	p.timeout = timeout
	p.pipeLineTask = f
	p.pipe = make(chan data[K, T])
	p.shutdown = make(chan struct{})
	p.wait = make(chan struct{})
	p.start()
	return p, nil
}
