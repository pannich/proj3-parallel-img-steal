// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	// "fmt"
	"image"
	"image/color"
)

// GetBoundary returns either the provided boundaries (specifically the first one if multiple are provided) or the bounds of the image itself.
func (img *Image) GetBoundary(boundaries ...image.Rectangle) image.Rectangle {
	// method for strucst *Image, take arguments and return image.Rectangle
	if len(boundaries) > 0 {
		return boundaries[0]
	} else {
		return img.Out.Bounds()
	}
}

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale(boundaries ...image.Rectangle) {
	bounds := img.GetBoundary(boundaries...)
	startX, endX, startY, endY := bounds.Min.X + 1, bounds.Max.X -1 , bounds.Min.Y, bounds.Max.Y

	// Adjustments for the boarders
	startY = Max(startY, 1)	// startY no lower than 0
	endY = Min(endY,(img.Bounds.Max.Y-1))	// endY no greater than original height index

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	// fmt.Printf("Bounds : %d %d %d %d\n", bounds.Min.Y, bounds.Max.Y, bounds.Min.X, bounds.Max.X)
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			var greyC uint16
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			// a os alpha : opacity
			//.RGBA() : Use these values for computations where you need standardized uint32 values,
			r, g, b, a := img.In.At(x, y).RGBA()

			// early return if pixel is black
			if checkPixelColor(r, g, b, a) {
				img.Out.Set(x, y, color.RGBA64{0, 0, 0, 0})
				continue
			} else {
				//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
				//For certain computations (i.e., convolution) the values might fall outside this
				// range so you need to clamp them between those values.
				// Create gray colour is to find 'average' of r g b
				greyC = clamp(float64(r+g+b) / 3)
				//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
				// for this that I won't get into right now).
				img.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}
}

//General convolution function
// pixel * kernel = newpixel
// kernel : https://www.songho.ca/dsp/convolution/convolution2d_example.html
func (img *Image) applyKernel(kernel []float64, bounds image.Rectangle) {

	startX, endX, startY, endY := bounds.Min.X + 1, bounds.Max.X - 1, bounds.Min.Y, bounds.Max.Y	// for bound [0:10] we do [1:9]

	// Adjustments for the boarders
	startY = Max(startY, 1)	// startY no lower than 1
	endY = Min(endY,(img.Bounds.Max.Y-1))	// endY no greater than original height index

	// fmt.Printf("Working Bounds : %d %d %d %d\n", startX, endX, startY, endY)
	// working on adjusted startY-endY
	for y := startY; y < endY; y++ {
			for x := startX; x < endX; x++ {	// width not including 1 pixel on the side
					var sumR, sumG, sumB, A float64
					// fmt.Printf("x,y : %d %d \n", x, y)

					r,g,b,a := img.In.At(x, y).RGBA()
					if checkPixelColor(r,g,b,a) {	// early exit if pixel i black
						img.Out.Set(x, y, color.RGBA64{0, 0, 0, 0})
						continue
					} else {
						for ky := -1; ky <= 1; ky++ {
								for kx := -1; kx <= 1; kx++ {
									boundaryClamp := func(val, min, max int) int {
										if val < min {
											return min
											} else if val > max {
												return max
											}
											return val
										}

										// Neighbours are not exceeding the full image boundary
										nx := boundaryClamp(x+kx, 0, bounds.Max.X-1)	// from index 0 to width-1
										ny := boundaryClamp(y+ky, 0, bounds.Max.Y+1)  // slice range

										r, g, b, a := img.In.At(nx, ny).RGBA() // get weight for a pixel by summing all the nb's weight; return uint32

										weight := kernel[(ky+1)*3+(kx+1)]			// get kernel weight from array
										sumR += weight * float64(r)
										sumG += weight * float64(g)
										sumB += weight * float64(b)
										A = float64(a)
								}
						}
						col := color.RGBA64{			// Ensure within bound
								R: uint16(clamp(sumR)),
								G: uint16(clamp(sumG)),
								B: uint16(clamp(sumB)),
								A: uint16(A),
						}

						// Save new pixel
						img.Out.Set(x, y, col)
						// fmt.Println(col)
					}
			}
	}

}

// Apply each kernel effect on pixel
func (img *Image) Sharpen(boundaries ...image.Rectangle) {
	kernel := []float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
	img.applyKernel(kernel, img.GetBoundary(boundaries...))
}

func (img *Image) EdgeDetection(boundaries ...image.Rectangle) {
	kernel := []float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
	img.applyKernel(kernel, img.GetBoundary(boundaries...))
}

func (img *Image) Blur(boundaries ...image.Rectangle) {
	kernel := []float64{
			1 / 9.0, 1 / 9.0, 1 / 9.0,
			1 / 9.0, 1 / 9.0, 1 / 9.0,
			1 / 9.0, 1 / 9.0, 1 / 9.0,
	}
	img.applyKernel(kernel, img.GetBoundary(boundaries...))
}

// --- Utils ----------------

func checkPixelColor(r,g,b,a uint32) bool {
	// Convert to non-pre-multiplied values and normalize to 0-255 range
	r, g, b, _ = r/257, g/257, b/257, a/257

	if r == 0 && g == 0 && b == 0 {
		return true
	} else {
		return false
	}
}

// --- End Utils ----------------
