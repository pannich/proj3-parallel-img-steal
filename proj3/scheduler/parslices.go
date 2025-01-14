package scheduler

import (
	"fmt"
	"math"
	"sync"
	"image"
	"time"

	"proj3/png"
)

// ---------------- Helper Function ---- //
// Sprawn go process by heights and wait till everyone's done.
// will slice by height. Each thread takes x rows to do task
func ProcessParallelSlices(pngImg *png.Image, numThreads int, effect string) (float64, *png.Image) {

	var wg sync.WaitGroup

	bounds := pngImg.Bounds					// have to access first then call Dx() on value it returns
	height := bounds.Dy()

	rowPerThread := int(math.Ceil(float64(height) / float64(numThreads)))

	actualNumThreads := png.Min(height, numThreads)  // min(no. all jobs, threadcount) // Ensure we do not spawn more threads than intervals

	// --------- start Parallel program for this image ---------
	start := time.Now()

	for i := 0; i < actualNumThreads; i++ {
		wg.Add(1)
		chunkStartY := rowPerThread * i
		chunkEndY := chunkStartY + rowPerThread
		if chunkEndY > height {
				chunkEndY = height
		}
		if chunkStartY < height {  // Only add work if the start is within the interval
				// fmt.Printf("Go %i is working\n", i)
				// note: pngImg and task is already a pointer

				// go ProcessImageChunk(pngImg, task, bounds.Min.X, chunkStartY, bounds.Max.X, chunkEndY, &wg)
				bounds := image.Rect(bounds.Min.X, chunkStartY, bounds.Max.X, chunkEndY)

				go func (bounds image.Rectangle) {
					// TODO
					defer wg.Done()		// will be called as soon as go routine completed

					switch effect {
					case "G" :
						pngImg.Grayscale(bounds)
					case "E":
						pngImg.EdgeDetection(bounds)
					case "S":
						pngImg.Sharpen(bounds)
					case "B":
						pngImg.Blur(bounds)
					}
				}(bounds)
		} else {
			wg.Done() // If no work is added, immediately call wg.Done
		}
	}
	wg.Wait() 	// done this effect
	// --------- End Parallel program for this effect ---------
	end := time.Since(start).Seconds()

	// fmt.Printf("  Parallelize Time for effect : %.2f\n" , end)		// Measure Parallelize time

	return end, pngImg
}

// ---------------- End Helper Function ---- //


// Main function pop each image from q
func RunParallelSlices(config Config) {

	// var pngImg *png.Image
	var time_ float64

	totalParallelTime := 0.0 	//accumulate parallel time

	queue, err := ReadTasksToQueue("../data/effects.txt", config.DataDirs)
	if err != nil {
		panic(err) // Or handle error more gracefully
	}

	// pop image from queue
	for {     // While.. .(until break)
		if queue.IsEmpty() {
				break
		}
		task := queue.Dequeue()

		pngImg, err := png.Load(task.InPath); if err != nil {panic(err)}

		// Performs a X filtering effect on the image
		for _, effect := range task.Effects {
			time_, pngImg = ProcessParallelSlices(pngImg, config.ThreadCount, effect)
			pngImg.In, pngImg.Out = pngImg.Out, pngImg.In		//Swap pointers
			totalParallelTime += time_
		}

		pngImg.Out, pngImg.In = pngImg.In, pngImg.Out
		//Saves the image to a new file
		err = pngImg.Save(task.OutPath); if err != nil {
			panic(err)		// check for error while saving
		}
		fmt.Printf("Accumulate Parallel Time 10 images: %.2f seconds\n", totalParallelTime)
	}
}
