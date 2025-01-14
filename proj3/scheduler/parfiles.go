package scheduler

import (
	"sync"
	"sync/atomic"
	"fmt"
	"time"

	"proj3/png"
)

type TASLock struct {
	state int32 // 0 indicates unlocked, 1 indicates locked
}

func (l *TASLock) Lock() {
	for atomic.CompareAndSwapInt32(&l.state, 0, 1) == false {
			// Loop until the state changes from 0 to 1
			// CompareAndSwapInt32(state, expected old, new if swap success)
			// Keep running if the comparison is wrong ie. state is 1 and expected old:0 (still locked)
	}
}

func (l *TASLock) Unlock() {
	atomic.StoreInt32(&l.state, 0)
}

// main function
func RunParallelFiles(config Config) {
	var wg sync.WaitGroup

	queue, err := ReadTasksToQueue("../data/effects.txt", config.DataDirs)
	if err != nil {
		panic(err) // Or handle error more gracefully
	}

	lock := &TASLock{}

	actualNumThreads := png.Min(queue.GetLength(), config.ThreadCount)

	start := time.Now()
	// --------- start Parallel program for 10 images ---------

	for i := 0; i < actualNumThreads; i++ {
			wg.Add(1)
			// fmt.Printf("Go %i is working\n", i)
			go func() {   // go routine starts anonymouse function
					for {     // While.. .(until break)
							lock.Lock()
							if queue.IsEmpty() {
									lock.Unlock()
									break
							}
							task := queue.Dequeue()
							lock.Unlock()

							ProcessImage(&task) // Your image processing function
					}
					wg.Done()
			}()
	}

	wg.Wait() // Wait for all goroutines to finish

	// --------- End Parallel program for 10 images ---------

	end := time.Since(start).Seconds()
	fmt.Printf("Parallelize Time : %.2f\n", end)		// Measure Parallelize time

}
