package go_pipeline

import "time"

type data struct {
	key string
	i   interface{}
}

type Pipeline struct {
	pipe         chan data
	pipeSize     int
	timeout      time.Duration
	pipeLineTask func(string, []interface{})
	shutdown     chan struct{}
	wait         chan struct{}
}

func (p *Pipeline) Add(key string, i interface{}) {
	p.pipe <- data{key, i}
}

func (p *Pipeline) clear(stateMap map[string][]interface{}) {
	for key, state := range stateMap {
		if len(state) > 0 {
			//fmt.Println("from timeout")
			p.pipeLineTask(key, state)
		}
	}
}

func (p *Pipeline) start() {
	go func() {
		var t *time.Timer
		stateMap := make(map[string][]interface{})
		forBreak := false
		for !forBreak {
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
				stateMap = make(map[string][]interface{})
			case <-p.shutdown:
				p.clear(stateMap)
				stateMap = make(map[string][]interface{})
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

func NewPipeline(pipeSize int, timeout time.Duration, f func(string, []interface{})) (*Pipeline, error) {
	p := &Pipeline{}
	p.pipeSize = pipeSize
	p.timeout = timeout
	p.pipeLineTask = f
	p.pipe = make(chan data)
	p.shutdown = make(chan struct{})
	p.wait = make(chan struct{})
	p.start()
	return p, nil
}
