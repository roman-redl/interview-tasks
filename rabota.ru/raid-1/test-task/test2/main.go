package main

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
)

const (
	inputDir       = 1
	outputDir      = 2
	fileNameSuffix = "_converted.jpeg"
)

func fileSize(file os.DirEntry) (int64, error) {
	fi, err := os.Stat(file.Name())
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func openReadFile(name string) (io.ReadCloser, error) {
	inputFile, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	log.Println("Load file " + name)
	return inputFile, nil
}

func openWriteFile(name string) (io.WriteCloser, error) {
	outputFile, err := os.Create(name + fileNameSuffix)
	if err != nil {
		return nil, err
	}
	log.Println("Create file " + name + fileNameSuffix)
	return outputFile, nil
}

func convert(im image.Image) *image.Gray {
	gray := image.NewGray(im.Bounds())
	bounds := im.Bounds()
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			gray.Set(x, y, im.At(x, y))
		}
	}
	return gray
}

func tryProcess(file os.DirEntry, inputDir, outputDir string) {
	if file.IsDir() {
		return
	}
	inputFile, err := openReadFile(inputDir + pathDelimiter + file.Name())
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = inputFile.Close()
	}()
	im, format, err := image.Decode(inputFile)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("found " + format + " image")

	outputFile, err := openWriteFile(outputDir + pathDelimiter + file.Name())
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = outputFile.Close()
	}()
	err = jpeg.Encode(outputFile, convert(im), &jpeg.Options{Quality: 100})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("No input/output dir")
	}
	files, err := os.ReadDir(os.Args[inputDir])
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		tryProcess(file, os.Args[inputDir], os.Args[outputDir])
	}
}
