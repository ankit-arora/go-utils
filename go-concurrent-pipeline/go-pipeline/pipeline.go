package go_pipeline

import "time"

type Pipeline struct {
	pipe         chan interface{}
	pipeSize     int
	timeout      time.Duration
	pipeLineTask func([]interface{})
	shutdown     chan struct{}
	wait         chan struct{}
}

func (p *Pipeline) Add(d interface{}) {
	p.pipe <- d
}

func (p *Pipeline) start() {
	go func() {
		var t *time.Timer
		var state []interface{}
		forBreak := false
		for !forBreak {
			t = time.NewTimer(p.timeout)
			select {
			case d := <-p.pipe:
				state = append(state, d)
				if len(state) == p.pipeSize {
					//fmt.Println("from data")
					p.pipeLineTask(state)
					state = nil
				}
			case <-t.C:
				if len(state) > 0 {
					//fmt.Println("from timeout")
					p.pipeLineTask(state)
					state = nil
				}
			case <-p.shutdown:
				if len(state) > 0 {
					//fmt.Println("from shutdown")
					p.pipeLineTask(state)
					state = nil
				}
				forBreak = true
			}
			t.Stop()
		}
		p.wait <- struct{}{}
	}()
}

func (p *Pipeline) Shutdown() {
	close(p.shutdown)
	<-p.wait
}

func NewPipeline(pipeSize int, timeout time.Duration, f func([]interface{})) (*Pipeline, error) {
	p := &Pipeline{}
	p.pipeSize = pipeSize
	p.timeout = timeout
	p.pipeLineTask = f
	p.pipe = make(chan interface{})
	p.shutdown = make(chan struct{})
	p.wait = make(chan struct{})
	p.start()
	return p, nil
}
