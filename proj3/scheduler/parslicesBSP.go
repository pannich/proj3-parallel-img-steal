package scheduler

import (
	"fmt"
	"image"
	"math"
	"sync"
	"time"

	"proj3/deque"
	"proj3/png"

	"math/rand"
)

type SharedContex struct {
	wgContext   *sync.WaitGroup
	mutex       *sync.Mutex
	cond        *sync.Cond
	counter     int
	threadCount int
}

func worker(id int, deques []*deque.DEQueue, ctx *SharedContex, img *png.Image, effect string) {
	myDeque := deques[id]
	wg := ctx.wgContext
	defer wg.Done()

	// image processing and work stealing loop
	for {
		var task *deque.Task
		var ok bool = true // Flag to check if task is obtained successfully, proceed to imageprocessing

		// Try to get a task from own deque first
		task, ok = myDeque.PopBottom()

		if !ok {
		// go to the stealing part if for some reason it fails ie. queue becomes empty
		attempted := make(map[int]bool)
		attempted[id] = true // Mark own deque as attempted
		for len(attempted) < len(deques) {
			i := rand.Intn(len(deques)) // Randomly select a deque index
			if _, found := attempted[i]; found || deques[i] == myDeque {
				continue // Skip if already attempted or is own deque
			}

			// No task in own deque, try to steal from others
				if !deques[i].IsEmpty() {
					fmt.Printf("Thread %d finished own deque %d, trying to steal from %d that has Bottom %d. And top %d\n", id, myDeque.GetBottom(), i, deques[i].GetBottom(), deques[i].GetTop())
					task, ok = deques[i].PopTop()
					if ok {
						// successfully stole a task
						break
					}
				}

				attempted[i] = true // Mark this deque as attempted and no task available
			}
		}

		if !ok {
			fmt.Printf("Go %d finished\n", id)
			break // No tasks available anywhere, exit
		}
		processImageSection(img, task.Bounds, effect)
	}

	/****barrier synchronization****/
	ctx.mutex.Lock()
	ctx.counter++
	if ctx.counter == ctx.threadCount {
		ctx.cond.Broadcast()
	} else {
		for ctx.counter != ctx.threadCount {
			ctx.cond.Wait()
		}
	}
	ctx.mutex.Unlock()
	/****barrier synchronization****/
}

func processImageSection(pngImg *png.Image, bounds image.Rectangle, effect string) {
	switch effect {
	case "G":
		pngImg.Grayscale(bounds)
	case "E":
		pngImg.EdgeDetection(bounds)
	case "S":
		pngImg.Sharpen(bounds)
	case "B":
		pngImg.Blur(bounds)
	}
}

// ---------------- Helper Function ---- //
// Sprawn go process by heights to apply one effect.
// Each thread takes x rows and wait until all threads are done.
// return time taken to process and the image with one effect applied.
func ProcessParallelSlicesBSP(pngImg *png.Image, numThreads int, effect string, optimized bool) (float64, *png.Image) {
	bounds := pngImg.Bounds
	height := bounds.Dy()

	/**** Adjust TaskCount per Thread here *****/
	taskCount := numThreads * 2

	chunkSize := int(math.Ceil(float64(height) / float64(taskCount)))
	actualNumThreads := png.Min(height, numThreads)

	deques := make([]*deque.DEQueue, actualNumThreads)
	for i := range deques {
		deques[i] = deque.NewDEQueue(taskCount/actualNumThreads + 1) // +1 in case of uneven division
	}
	// Calculate the number of tasks per deque
	tasksPerDeque := int(math.Ceil(float64(taskCount) / float64(actualNumThreads)))

	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx := SharedContex{wgContext: &wg, mutex: &mu, cond: sync.NewCond(&mu), counter: 0, threadCount: actualNumThreads}

	start := time.Now()

	//--- Optimized Enque version ---
	if optimized {
		// Parallel Distribute tasks across deques in blocks by locality
		var wg2 sync.WaitGroup
		for i := 0; i < actualNumThreads; i++ {
			wg2.Add(1)
			go func(workerId int) {
				defer wg2.Done()
				SequentialEnque(workerId, actualNumThreads, tasksPerDeque, taskCount, chunkSize, bounds, deques)
				}(i)
			}
			wg2.Wait() // Synchrnoize the task enqueue
			} else {
				// Sequential Distribute
				for i := 0; i < actualNumThreads; i++ {
					SequentialEnque(i, actualNumThreads, tasksPerDeque, taskCount, chunkSize, bounds, deques)
				}
			}

	// Start all workers after enque
	//each worker can see the everyone deques and shared context
	for i := 0; i < actualNumThreads; i++ {
		ctx.wgContext.Add(1)
		go worker(i, deques, &ctx, pngImg, effect)
	}

	// Synchronize again to catch all the image processing go routines.
	ctx.wgContext.Wait()

	end := time.Since(start).Seconds()

	return end, pngImg
}

// ---------------- End Helper Function ---- //

// Main function pop each image from q
func RunParallelSlicesBSP(config Config, optimized bool) {

	// var pngImg *png.Image
	var time_ float64

	totalParallelTime := 0.0 //accumulate parallel time

	queue, err := ReadTasksToQueue("../data/effects.txt", config.DataDirs)
	if err != nil {
		panic(err) // Or handle error more gracefully
	}

	// pop image from queue
	// While.. .(until break)
	for {
		if queue.IsEmpty() {
			break
		}
		task := queue.Dequeue()

		pngImg, err := png.Load(task.InPath)
		if err != nil {
			panic(err)
		}

		// Performs an effect on the image
		for _, effect := range task.Effects {
			time_, pngImg = ProcessParallelSlicesBSP(pngImg, config.ThreadCount, effect, optimized)
			pngImg.In, pngImg.Out = pngImg.Out, pngImg.In //Swap pointers
			totalParallelTime += time_
		}

		pngImg.Out, pngImg.In = pngImg.In, pngImg.Out
		//Saves the image to a new file
		err = pngImg.Save(task.OutPath)
		if err != nil {
			panic(err) // check for error while saving
		}
		fmt.Printf("Accumulate Parallel Time 10 images: %.2f seconds\n", totalParallelTime)
	}
}
