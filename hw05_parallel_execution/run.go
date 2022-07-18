package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {

	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	errorsCh := make(chan error, len(tasks))
	completedCh := make(chan struct{})

	defer func() {
		close(completedCh)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go doTask(tasksCh, errorsCh, completedCh, &wg)
	}

	for _, task := range tasks {
		tasksCh <- task

		if m > 0 && len(errorsCh) >= m {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}

func doTask(tasks chan Task, errors chan error, completed chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-completed:
			wg.Done()
			return
		case task := <-tasks:
			if err := task(); err != nil {
				errors <- err
			}
		}
	}
}
