package main

import (
	"io/fs"
	"log"
	"os"
	"sync"
)

// TODO change go func to worker pool
func main() {
	if len(os.Args) < 3 {
		log.Fatal("No input/output dir")
	}
	files, err := os.ReadDir(os.Args[inputDir])
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan struct{}, 50)
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go func(file fs.DirEntry) {
			ch <- struct{}{}
			defer wg.Done()
			tryProcess(file, os.Args[inputDir], os.Args[outputDir])
			<-ch
		}(file)
	}

	wg.Wait()
}
