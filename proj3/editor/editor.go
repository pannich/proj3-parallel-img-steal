/*
editor.go takes in 3 arguments
which will be initialize in config : input files, mode, threadcount
then pass 'config' to the scheduler.Schedule(config)
which apply the effects
*/

package main

import (
	"fmt"		// for formatted I/O operations
	"os"
	"proj3/scheduler"		// scheduling and managing the image processing tasks
	"strconv"
	"time"
)


const usage = "Usage: editor data_dir mode [number of threads]\n" +
	"data_dir = The data directory to use to load the images.\n" +
	"mode     = (s) run sequentially, (parfiles) process multiple files in parallel, (parslices) process slices of each image in parallel \n" +
	"[number of threads] = Runs the parallel version of the program with the specified number of threads.\n"

func main() {

	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	// set up config. Config struct form scheduler package.
	// - DataDirs : Directory from which images should be loaded -- First command line arg
	// - Mode : 's', 'parfiles', 'parslices'
	// - ThreadCount : if not provide is Sequential mode
	// ie $: go run editor.go big+small pipeline 2

	config := scheduler.Config{DataDirs: "", Mode: "", ThreadCount: 0}
	config.DataDirs = os.Args[1]

	if len(os.Args) >= 3 {
		config.Mode = os.Args[2]
		threads, _ := strconv.Atoi(os.Args[3])
		config.ThreadCount = threads
	} else {
		config.Mode = "s"
	}
	start := time.Now()
	scheduler.Schedule(config)		// run task
	end := time.Since(start).Seconds()
	fmt.Printf("  Program time: %.2f\n", end)		// Measure Program time

}
