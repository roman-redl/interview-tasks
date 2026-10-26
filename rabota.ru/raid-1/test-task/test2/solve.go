package main

// Решение задачи ImageConverter (rabota.ru, тест 2) — закрытие TODO из
// черновика решения: worker pool вместо горутины на файл + лимит памяти.
//
// 1. Распараллеливание: пул воркеров фиксированного размера (CPU-bound
//    конвертация — больше NumCPU*2 воркеров бессмысленно), задачи —
//    через канал, а не go func на файл: тысячи картинок не создают
//    тысяч горутин.
// 2. Не падать при большом количестве: никакой ошибки из обработки
//    отдельного файла не fatal — лог и следующая; память ограничена п.3.
// 3. Не более X МБ картинок в памяти одновременно: взвешенный семафор
//    с зерном 1 МБ (ceil(size/1MB), минимум 1): воркер занимает столько
//    зерён, сколько весит файл, перед чтением и освобождает после записи.
//    Третий аргумент командной строки — X; фиксированного числа
//    одновременно обрабатываемых картинок больше нет: их определяет вес.
//
// Файл самодостаточен: собирается и запускается `go run solve.go`
// без main.go (поэтому здесь свои копии хелперов условия и filepath.Join
// вместо константы pathDelimiter с build-тегами).

import (
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	inputDirArg  = 1
	outputDirArg = 2
	mbLimitArg   = 3
	fileSuffix   = "_converted.jpeg"
	mb           = 1024 * 1024
)

// memLimiter — взвешенный семафор: зерна по 1 МБ, емкость X МБ.
type memLimiter struct {
	grains chan struct{}
}

func newMemLimiter(megabytes int) *memLimiter {
	if megabytes < 1 {
		megabytes = 1
	}
	return &memLimiter{grains: make(chan struct{}, megabytes)}
}

func (l *memLimiter) weight(size int64) int {
	w := int((size + mb - 1) / mb) // ceil, min 1
	if w < 1 {
		w = 1
	}
	if w > cap(l.grains) { // файл больше лимита — всё равно обработаем, заняв весь
		w = cap(l.grains)
	}
	return w
}

func (l *memLimiter) acquire(w int) {
	for i := 0; i < w; i++ {
		l.grains <- struct{}{}
	}
}

func (l *memLimiter) release(w int) {
	for i := 0; i < w; i++ {
		<-l.grains
	}
}

type job struct {
	path   string
	size   int64
	weight int
}

func openReadFile(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	log.Println("Load file " + name)
	return f, nil
}

func openWriteFile(name string) (io.WriteCloser, error) {
	f, err := os.Create(name + fileSuffix)
	if err != nil {
		return nil, err
	}
	log.Println("Create file " + name + fileSuffix)
	return f, nil
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

func tryProcess(j job, lim *memLimiter) {
	lim.acquire(j.weight)
	defer lim.release(j.weight)

	inputFile, err := openReadFile(j.path)
	if err != nil {
		log.Println(err)
		return
	}
	defer inputFile.Close()

	im, format, err := image.Decode(inputFile)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("found " + format + " image")

	outPath := filepath.Join(os.Args[outputDirArg], filepath.Base(j.path))
	outputFile, err := openWriteFile(outPath)
	if err != nil {
		log.Println(err)
		return
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, convert(im), &jpeg.Options{Quality: 100}); err != nil {
		log.Println(err)
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("No input/output dir")
	}
	mbLimit := runtime.NumCPU() * 8 // дефолт, если X не передан
	if len(os.Args) > mbLimitArg {
		if n, err := fmt.Sscanf(os.Args[mbLimitArg], "%d", &mbLimit); n != 1 || err != nil || mbLimit < 1 {
			log.Fatalf("bad MB limit %q", os.Args[mbLimitArg])
		}
	}

	files, err := os.ReadDir(os.Args[inputDirArg])
	if err != nil {
		log.Fatal(err)
	}

	lim := newMemLimiter(mbLimit)
	jobs := make(chan job)

	var wg sync.WaitGroup
	workers := runtime.NumCPU() * 2
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				tryProcess(j, lim)
			}
		}()
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(os.Args[inputDirArg], f.Name())
		fi, err := os.Stat(path)
		if err != nil {
			log.Println(err)
			continue
		}
		jobs <- job{path: path, size: fi.Size(), weight: lim.weight(fi.Size())}
	}
	close(jobs)
	wg.Wait()
}
