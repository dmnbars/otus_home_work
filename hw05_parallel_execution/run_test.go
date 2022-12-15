package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("test for concurrent without sleeps by sort", func(t *testing.T) {
		tasksCount := 1000
		tasks := make([]Task, 0, tasksCount)
		ch := make(chan int, tasksCount)
		processed := make([]int, 0, tasksCount)

		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			i := i
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				ch <- i
				return nil
			})
		}

		err := Run(tasks, 5, 1)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

		for val := range ch {
			processed = append(processed, val)
			if len(processed) == tasksCount {
				break
			}
		}

		require.False(t, sort.IsSorted(sort.IntSlice(processed)), "tasks were run sequentially?")
	})

	t.Run("tasks count less when workers", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, 10, 1)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("negative count of workers", func(t *testing.T) {
		err := Run([]Task{}, -1, 5)
		require.ErrorIs(t, err, ErrNotPositiveCountOfWorkers)
	})

	t.Run("zero count of workers", func(t *testing.T) {
		err := Run([]Task{}, 0, 5)
		require.ErrorIs(t, err, ErrNotPositiveCountOfWorkers)
	})

	t.Run("negative errors limit", func(t *testing.T) {
		err := Run([]Task{}, 1, -1)
		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
	})

	t.Run("zero errors limit", func(t *testing.T) {
		err := Run([]Task{}, 1, 0)
		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
	})
}
