package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, 0)
	errCh := make(chan error, 0)
	quit := make(chan struct{}, 0)
	resultCh := make(chan error, 0)

	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		go doWork(tasksCh, errCh, wg)
	}

	go sendTasks(tasks, tasksCh, quit, resultCh)
	go checkErrors(errCh, m, quit)

	result := <-resultCh

	wg.Wait()

	close(quit)
	close(errCh)
	close(resultCh)

	return result
}

func sendTasks(tasks []Task, tasksCh chan<- Task, quit <-chan struct{}, resultCh chan<- error) {
	defer close(tasksCh)
	for _, task := range tasks {
		select {
		case <-quit:
			resultCh <- ErrErrorsLimitExceeded
			return
		case tasksCh <- task:
		}
	}

	resultCh <- nil
}

func checkErrors(errCh <-chan error, m int, quit chan<- struct{}) {
	for err := range errCh {
		if err != nil {
			m--
		}
		if m == 0 {
			quit <- struct{}{}
		}
	}
}

func doWork(tasksCh <-chan Task, errCh chan<- error, wg *sync.WaitGroup) {
	wg.Add(1)

	defer func() {
		wg.Done()
	}()

	for task := range tasksCh {
		errCh <- task()
	}
}
