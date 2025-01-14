// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"fmt"
)

// The Image represents a structure for working with PNG images.
// You are allowed to update this and change it as you wish!
type Image struct {
	In     *image.RGBA64   //The original pixels before applying the effect
	Out    *image.RGBA64   //The updated pixels after applying teh effect
	Bounds image.Rectangle //The size of the image. By default is full image
}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
// You are allowed to modify and update this as you wish
func Load(filePath string) (*Image, error) {
	// output returns pointer to an Image structure, and an error

	inReader, err := os.Open(filePath) 		// If file doesn't exist, returns error

	if err != nil {				// Check if can open file
		return nil, err
	}
	defer inReader.Close()

	inOrig, err := png.Decode(inReader)		// package image/png. decode te image as .png
																				// if success inOrig is in image.Image format

	if err != nil {				// decoding fails
		return nil, err
	}

	bounds := inOrig.Bounds() 	// get size of original image. This will be used to set size of new image.

	outImg := image.NewRGBA64(bounds)		// create new images outImg and inImg which are initially blank. // type : image.RGBA64
	inImg := image.NewRGBA64(bounds)

	// loop thru inOrig and copy into the inImg. We will use inImg for image processing.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := inOrig.At(x, y).RGBA()
			inImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}

	// initialize Image struct
	task := &Image{}
	task.In = inImg				// copied image
	task.Out = outImg			// blank image
	task.Bounds = bounds	// image boundary
	return task, nil
}

// Save saves the image to the given file
// You are allowed to modify and update this as you wish
func (img *Image) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.Out)
	if err != nil {
		return err
	}
	return nil
}


// To compare result with expected
func CompareImages(img1, img2 *Image) bool {
	if img1.Bounds != img2.Bounds {
			fmt.Println("Images have different dimensions")
			return false
	}
	areIdentical := true
	for y := 0; y < img1.Bounds.Dy(); y++ {
			for x := 0; x < img1.Bounds.Dx(); x++ {
					if img1.In.At(x, y) != img2.Out.At(x, y) {
						fmt.Printf("Pixel difference at %d,%d: %v != %v\n", x, y, img1.In.At(x, y), img2.Out.At(x, y))
						areIdentical = false
					}
			}
	}
	return areIdentical
}

func (img *Image) SetInput(startX int, endX int, startY int, endY int){
	// copy to input

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			r, g, b, a := img.Out.At(x, y).RGBA()
			img.In.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
}

//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}
