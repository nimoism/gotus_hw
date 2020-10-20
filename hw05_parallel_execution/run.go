package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

type Result struct {
	err error
}

func (r Result) Err() error {
	return r.err
}

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrNoWorkers = errors.New("no workers to handle tasks")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, workerCount int, maxErrorCount int) error {
	if workerCount < 1 {
		return ErrNoWorkers
	}
	ignoreErrors := maxErrorCount < 1
	taskCh := make(chan Task)
	resCh := make(chan Result)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		runWorkers(workerCount, taskCh, resCh, doneCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		pushTasks(taskCh, tasks, doneCh)
	}()

	errorCount := 0
	var err error

	for {
		result, ok := <-resCh
		if !ok {
			break
		}
		if ignoreErrors {
			continue
		}
		if result.Err() != nil {
			errorCount++
		}
		if errorCount >= maxErrorCount {
			err = ErrErrorsLimitExceeded
			close(doneCh)
			break
		}
	}
	wg.Wait()
	return err
}

func pushTasks(taskCh chan<- Task, tasks []Task, doneCh <-chan struct{}) {
	defer close(taskCh)
	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-doneCh:
			return
		}
	}
}

func runWorkers(count int, taskCh <-chan Task, resCh chan<- Result, doneCh <-chan struct{}) {
	defer close(resCh)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			runWorker(taskCh, resCh, doneCh)
		}()
	}
	wg.Wait()
}

func runWorker(taskCh <-chan Task, resCh chan<- Result, doneCh <-chan struct{}) {
	for {
		task, ok := <-taskCh
		if !ok {
			return
		}
		err := task()
		select {
		case resCh <- Result{err: err}:
		case <-doneCh:
			return
		}
	}
}
