package go_shutdown_hook

import (
	"os"
	"os/signal"
	"syscall"
)

var shutdownSignals chan os.Signal
var funcChannel chan func()
var isStarted bool = false
var funcArray []func()
var done chan struct{}

func ADD(f func()) {
	if !isStarted {
		shutdownSignals = make(chan os.Signal, 1)
		funcChannel = make(chan func())
		done = make(chan struct{})
		signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)
		start()
		isStarted = true
	}
	funcChannel <- f
}

func Wait() {
	<-done
}

func executeHooks() {
	for _, f := range funcArray {
		f()
	}
}

func start() {
	go func() {
		shutdown := false
		for !shutdown {
			select {
			case <-shutdownSignals:
				executeHooks()
				shutdown = true
			case f := <-funcChannel:
				funcArray = append(funcArray, f)
			}
		}
		done <- struct{}{}
	}()
}
