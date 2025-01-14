package scheduler

import (
	"proj3/png"
	// "fmt"
)

// Take ImageTask that contains all image info and effects
func ProcessImage(task *ImageTask) {
	// Load from path to return *Image
	pngImg, err := png.Load(task.InPath)
	if err != nil {
		panic(err)
	}

	// Performs a X filtering effect on the image
	for _, effect := range task.Effects {
		// apply effect
		switch effect {
			case "G" :
				pngImg.Grayscale()		// object from the same package call method in the same package no need to do png.Grayscale().
															// note : The method call does not require the package name as a prefix; it only requires a valid instance of the type.
			case "E":
				pngImg.EdgeDetection()
			case "S":
				pngImg.Sharpen()
			case "B":
				pngImg.Blur()
		}
		pngImg.In, pngImg.Out = pngImg.Out, pngImg.In
		}
		//Saves the image to a new file
		pngImg.Out, pngImg.In = pngImg.In, pngImg.Out
		err = pngImg.Save(task.OutPath); if err != nil {
			panic(err)		// check for error while saving
		}
}
