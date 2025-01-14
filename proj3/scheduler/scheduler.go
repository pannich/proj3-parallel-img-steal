package scheduler

type Config struct {
	DataDirs string //Represents the data directories to use to load the images.
	Mode     string // Represents which scheduler scheme to use
	ThreadCount int // Runs parallel version with the specified number of threads
}

// ImageTask details from effects.txt
type ImageTask struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
}

//Run the correct version based on the Mode field of the configuration value
func Schedule(config Config) {
	if config.Mode == "s" {
		RunSequential(config)
	} else if config.Mode == "parfiles" {
		RunParallelFiles(config)
	} else if config.Mode == "parslices" {
		RunParallelSlices(config)
	} else if config.Mode == "parslicesBSP" {
		RunParallelSlicesBSP(config, false)
	} else if config.Mode == "parslicesBSPOptimized" {
		RunParallelSlicesBSP(config, true)
	} else {
		panic("Invalid scheduling scheme given.")
	}
}
