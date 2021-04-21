package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"math/bits"
	"os"

	"github.com/nfnt/resize"
)

func main() {

	filename1Ptr := flag.String("file1", "", "filepath to generate hash for")
	filename2Ptr := flag.String("file2", "", "filepath to generate hash for")
	flag.Parse()

	hash1 := openAndHash(*filename1Ptr)
	hash2 := openAndHash(*filename2Ptr)

	fmt.Printf("file1: %064b, %016X\n", hash1, hash1)
	fmt.Printf("file2: %064b, %016X\n", hash2, hash2)

	difference := hash1 ^ hash2
	fmt.Println(bits.OnesCount64(uint64(difference)))

	// hashedImage := hashedImaged(hash1)
	// saveImg("hashed.png", hashedImage)
}

func hashToByteArray(hash uint64) []uint8 {

	colors := make([]uint8, 0, 64)

	for y := 7; y >= 0; y-- {

		offset := y * 8
		row := (hash & (0xFF << (offset))) >> offset

		for x := 7; x >= 0; x-- {
			pxl := uint8(255)
			if (row & (1 << x)) > 0 {
				pxl = 0
			}

			colors = append(colors, pxl)
		}

	}
	return colors
}

func hashedImaged(hash uint64) *image.Gray {
	byteArr := hashToByteArray(hash)

	img := image.NewGray(image.Rect(0, 0, 8, 8))

	img.Pix = byteArr

	return img
}

func openAndHash(filepath string) uint64 {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	result := AHash(img)

	return result
}

func imgToGrayscale(img image.Image) *image.Gray {
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}
	return grayImg
}

func imgResize(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.NearestNeighbor)
}

func AHash(img image.Image) uint64 {
	if img.Bounds().Dx() != 8 && img.Bounds().Dy() != 8 {
		// fmt.Println("Resizing to 8x8")
		resized := imgResize(img, 8, 8)
		img = resized
	}

	// saveImg("resized.png", img)

	gray := imgToGrayscale(img)
	// saveImg("gray.png", gray)

	avg := average(gray.Pix)

	return AHashComputeBits(gray.Pix, avg)
}

func AHashComputeBits(data []uint8, avg uint8) uint64 {
	hash := uint64(0)

	for idx, pix := range data {
		if idx > 0 {
			hash <<= 1
		}

		if pix >= avg {
			hash |= 0b1
		}
		// fmt.Printf("%064b, %03d : %03d\n", hash, avg, pix)
	}

	return hash
}

func average(arr []uint8) uint8 {
	total := uint32(0)

	if len(arr) == 0 {
		return 0
	}
	for _, el := range arr {
		total += uint32(el)
	}

	avg := float64(total) / float64(len(arr))

	return uint8(math.Round(avg))
}

func saveImg(filename string, img image.Image) {

	file, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}
