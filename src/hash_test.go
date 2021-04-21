package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"testing"
)

var (
	testFileSmall         = "../testdata/shion.png"
	testWidthResize  uint = 32
	testHeightResize uint = 32
)

func loadFile(filepath string) image.Image {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return img
}

func BenchmarkGrayscaleResizeSmall(b *testing.B) {
	img := loadFile(testFileSmall)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gray := imgToGrayscale(img)
		imgResize(gray, testWidthResize, testHeightResize)
	}
}

func BenchmarkResizeGrayscaleSmall(b *testing.B) {
	img := loadFile(testFileSmall)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resized := imgResize(img, testWidthResize, testHeightResize)
		imgToGrayscale(resized)
	}
}

func TestAverageSimple(t *testing.T) {
	testCases := []struct {
		arr []uint8
		avg uint8
	}{
		{},
		{
			arr: []uint8{1, 2, 3},
			avg: 2,
		},
		{
			arr: []uint8{1, 2, 3, 4, 5},
			avg: 3,
		}, {
			arr: []uint8{255, 255, 255},
			avg: 255,
		}, {
			arr: []uint8{3, 2},
			avg: 3, // round up
		},
	}

	for _, testCase := range testCases {
		arr := testCase.arr
		expected := testCase.avg

		actual := average(arr)
		if actual != expected {
			t.Errorf("got %d; want %d", actual, expected)
		}
	}

}

func TestAHashComputeBits(t *testing.T) {
	testCases := []struct {
		arr  []uint8
		avg  uint8
		hash uint64
	}{
		{},
		{
			arr:  []uint8{1, 2, 3},
			avg:  2,
			hash: 0b011,
		},
		{
			arr:  []uint8{1, 2, 3, 4, 5},
			avg:  3,
			hash: 0b00111,
		},
		{
			arr:  []uint8{1, 2, 3, 4, 5},
			avg:  0,
			hash: 0b11111,
		},
		{
			arr:  []uint8{1, 2, 3, 4, 5},
			avg:  0,
			hash: 0b11111,
		},
		{
			arr:  []uint8{1, 2, 3, 4, 5},
			avg:  10,
			hash: 0,
		},
	}

	for _, testCase := range testCases {
		arr := testCase.arr
		expected := testCase.hash

		actual := AHashComputeBits(arr, testCase.avg)
		if actual != expected {
			t.Errorf("got %d; want %d", actual, expected)
		}
	}
}

func TestHashToByteArray(t *testing.T) {
	var num uint64 = 0x0badfeeddeadbeef

	result := hashToByteArray(num)
	fmt.Printf("%016b\n", 0x0bad)

	fmt.Printf("%v\n", result)
	fmt.Printf("%d\n", len(result))

}
