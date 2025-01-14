/*
sequential.go
schedule task sequentially
put images in arraylist and execute one by one
*/

package scheduler

import (
	"encoding/json"
	"os"

	"strings"
	"bufio"
)

// ReadImageTasks reads and decodes the JSON tasks from effects.txt; then put in slice (similar to arraylist)
// put the object in effects.txt into the slice
func ReadImageTasks(effectsPathFile string) ([]ImageTask, error) {
	effectsFile, err := os.Open(effectsPathFile)   // open file readonly mode
	if err != nil {
			return nil, err
	}
	defer effectsFile.Close()

	var tasks []ImageTask

	scanner := bufio.NewScanner(effectsFile)			// Read line by line
	for scanner.Scan() {
    var task ImageTask
    err := json.Unmarshal(scanner.Bytes(), &task)			// Unmarshal takes in []byte slice and copy into &task struct
    if err != nil {
        return nil, err  // Handle the error appropriately
    }
    tasks = append(tasks, task)
}
	return tasks, nil
}


//Main Function
// take in config which is userinput
func RunSequential(config Config) {
	// Load and put tasks in slice
	tasks, err := ReadImageTasks("../data/effects.txt")		// return tasks slice
	if err != nil {
		panic(err) // Or handle error more gracefully
	}

	// Assuming DataDirs might contain multiple directories separated by '+'
	dataDirs := strings.Split(config.DataDirs, "+")

	// Process each directory
	for _, data_dir := range dataDirs {
		for _, task := range tasks {

			SetTaskPath(&task, data_dir)

			ProcessImage(&task)
		}
	}
}
