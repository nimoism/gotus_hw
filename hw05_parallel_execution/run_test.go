package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
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
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
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
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("no tasks does nothing", func(t *testing.T) {
		tasks := make([]Task, 0)

		result := Run(tasks, 10, 10)
		require.Nil(t, result)
	})

	t.Run("no workers error", func(t *testing.T) {
		taskCount := 10
		for _, tst := range []struct {
			workerCount int
		}{
			{workerCount: 0},
			{workerCount: -1},
		} {
			t.Run(fmt.Sprintf("worker count %d", tst.workerCount), func(t *testing.T) {
				tasks := make([]Task, 0, taskCount)
				var runTasksCount int32
				for i := 0; i < taskCount; i++ {
					taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
					tasks = append(tasks, func() error {
						time.Sleep(taskSleep)
						atomic.AddInt32(&runTasksCount, 1)
						return nil
					})
				}
				result := Run(tasks, tst.workerCount, 0)
				require.Equal(t, ErrNoWorkers, result)
				require.Equal(t, int32(0), runTasksCount)
			})
		}
	})

	t.Run("no max error count", func(t *testing.T) {
		taskCount := 20
		workerCount := 10
		for _, tst := range []struct {
			maxErrorCount int
		}{
			{maxErrorCount: 0},
			{maxErrorCount: -1},
		} {
			t.Run(fmt.Sprintf("max error count %d", tst.maxErrorCount), func(t *testing.T) {
				tasks := make([]Task, 0, taskCount)
				var runTasksCount int32
				for i := 0; i < taskCount; i++ {
					taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

					tasks = append(tasks, func() error {
						time.Sleep(taskSleep)
						atomic.AddInt32(&runTasksCount, 1)
						return fmt.Errorf("error from task %d", i)
					})
				}
				result := Run(tasks, workerCount, tst.maxErrorCount)
				require.NoError(t, result)
				require.Equal(t, int32(taskCount), runTasksCount)
			})
		}
	})

	t.Run("worker count more than task count", func(t *testing.T) {
		taskCount := 10
		workerCount := 20
		for _, mecTst := range []struct {
			mec int
			err error
		}{
			{mec: workerCount + 1},
			{mec: workerCount},
			{mec: workerCount - 1},
			{mec: taskCount + 1},
			{mec: taskCount, err: ErrErrorsLimitExceeded},
			{mec: taskCount - 1, err: ErrErrorsLimitExceeded},
			{mec: 0},
			{mec: -1},
		} {
			t.Run(fmt.Sprintf("max error count %d", mecTst.mec), func(t *testing.T) {
				tasks := make([]Task, 0, taskCount)
				var runTasksCount int32
				for i := 0; i < taskCount; i++ {
					taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

					tasks = append(tasks, func() error {
						time.Sleep(taskSleep)
						atomic.AddInt32(&runTasksCount, 1)
						return fmt.Errorf("error from task %d", i)
					})
				}
				result := Run(tasks, workerCount, mecTst.mec)
				require.Equal(t, mecTst.err, result)
				require.Equal(t, int32(taskCount), runTasksCount)
			})
		}
	})

}
