package goroutines_pool

import (
	"errors"
	"sync"
)

type Pool struct {
	funcChannel chan func()
	semaphore   chan struct{}
	stopped     chan struct{}
}

func (p *Pool) Submit(f func()) {
	p.funcChannel <- f
}

func (p *Pool) start() {
	go func() {
		var done sync.WaitGroup
		for f := range p.funcChannel {
			p.semaphore <- struct{}{}
			done.Add(1)
			go func(f func()) {
				f()
				done.Done()
				<-p.semaphore
			}(f)
		}
		done.Wait()
		close(p.semaphore)
		p.stopped <- struct{}{}
	}()
}

func (p *Pool) Shutdown() {
	close(p.funcChannel)
	<-p.stopped
}

func NewPool(maxConcurrency int) (*Pool, error) {
	if maxConcurrency <= 0 {
		return nil, errors.New("maxConcurrency cannot be 0 or negative")
	}
	p := &Pool{}
	p.funcChannel = make(chan func())
	p.stopped = make(chan struct{})
	p.semaphore = make(chan struct{}, maxConcurrency)
	p.start()
	return p, nil
}
