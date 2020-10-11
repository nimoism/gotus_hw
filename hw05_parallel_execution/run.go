package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

// Notifier allows notify once and ignore other notifications.
type Notifier struct {
	C chan struct{}
}

func (n Notifier) Notify() {
	select {
	case n.C <- struct{}{}:
	default:
	}
}

func NewNotifier() Notifier {
	return Notifier{C: make(chan struct{}, 1)}
}

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
	stop := NewNotifier()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		runWorkers(workerCount, taskCh, resCh, doneCh)
		stop.Notify()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		pushTasks(taskCh, tasks, doneCh)
		wg.Done()
	}()

	errorCount := 0
	var err error
	for {
		select {
		case result := <-resCh:
			if ignoreErrors {
				continue
			}
			if result.Err() != nil {
				errorCount++
			}
			if errorCount >= maxErrorCount {
				err = ErrErrorsLimitExceeded
				stop.Notify()
			}
		case <-stop.C:
			close(doneCh)
			wg.Wait()
			return err
		}
	}
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
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			runWorker(taskCh, resCh, doneCh)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runWorker(taskCh <-chan Task, resCh chan<- Result, doneCh <-chan struct{}) {
	for {
		select {
		case task, ok := <-taskCh:
			if !ok {
				return
			}
			err := task()
			select {
			case resCh <- Result{err: err}:
			case <-doneCh:
				return
			}
		case <-doneCh:
			return
		}
	}
}
