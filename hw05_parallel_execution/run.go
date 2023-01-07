package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded    = errors.New("errors limit exceeded")
	ErrErrorsGoroutinesNumber = errors.New("negative or zero number of goroutines ")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrErrorsGoroutinesNumber
	}

	taskChan := make(chan Task, n)
	errChan := make(chan struct{}, m)
	cancel := make(chan struct{}, 1)

	var hasError error

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

		TASK:
			for task := range taskChan {
				select {
				case <-cancel:
					break TASK
				default:
					if err := task(); err != nil {
						errChan <- struct{}{}
					}
				}
			}
		}()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		errCount := 0

		for range errChan {
			errCount++

			if errCount == m {
				close(cancel)
				break
			}
		}
	}()

LOOP:
	for _, task := range tasks {
		select {
		case <-cancel:
			hasError = ErrErrorsLimitExceeded
			break LOOP
		default:
			taskChan <- task
		}
	}

	close(taskChan)

	if hasError == nil {
		close(errChan)
	}

	wg.Wait()

	return hasError
}
