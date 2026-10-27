# Rabota.ru — заход 1, тестовое задание

- Даты: 2022
- Формат: домашнее задание, два задания на Go — исходный код проекта в условии, доработать его

## 1. QueueAPI: битые сообщения и деградация под нагрузкой

Сервис QueueAPI на fasthttp принимает PUT-запросы и асинхронно кладёт тело запроса в очередь сообщений (сетевая отправка в примере заменена на задержку и лог). Две проблемы: при большом трафике и больших входящих сообщениях сообщения в очереди получаются сломанными — впечатление, что они состоят из кусков разных сообщений; а при большом трафике сервис съедает все ресурсы и умирает. Нужно, чтобы он пусть и тормозил, но выживал, при этом затормаживая и потребителей: потребители сами должны решать, готовы ли подождать подольше перед отправкой в очередь.

```go
package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func putToQueue(body []byte) {
	//тут по сети кладем какие-то байты в очередь
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
	log.Println(body)
}

func main() {
	r := router.New()

	r.PUT("/queue", func(ctx *fasthttp.RequestCtx) {
		go func(body []byte) {
			putToQueue(body)
		}(ctx.Request.Body())
		ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
	})

	err := fasthttp.ListenAndServe(":8088", r.Handler)
	log.Fatal(err)
}
```

Зависимости: fasthttp, fasthttp/router, envviper (go.mod прилагается).

## 2. ImageConverter: параллельная конвертация картинок

Утилита читает картинки (gif/jpeg/png) из входной директории, конвертирует в оттенки серого и сохраняет JPEG с суффиксом `_converted.jpeg` в выходную директорию; пути передаются первыми двумя аргументами командной строки. Требуется: распараллелить обработку картинок; при слишком большом их количестве — не падать; в идеале — считывать в память не более X Мб картинок, где X — третий аргумент командной строки, вместо фиксированного числа одновременно обрабатываемых картинок.

```go
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
```

Разделитель пути задаётся константой `pathDelimiter` через build-теги: `linux.go` (`/`, `//go:build !windows`) и `win.go` (`\`, `//go:build windows`); в main.go также есть неиспользуемая функция `fileSize`.

## Решения

Папки solution/ нет: решения — сами проекты. В каждом `main.go` — исходный код из условия, `solve.go` — решение кандидата.

- Задача 1 — `test1/`: `solve.go` — копирование тела запроса (буфер fasthttp переиспользуется — источник «битых» сообщений), пул воркеров + очередь конечной глубины, при переполнении — 429 с `Retry-After` и глубиной очереди в заголовке (backpressure). Сборка и запуск: `go run solve.go` (воркеры/глубина — env `QUEUEAPI_WORKERS`, `QUEUEAPI_BACKLOG`).
- Задача 2 — `test2/`: `solve.go` — worker pool вместо горутины на файл; лимит памяти X МБ из третьего аргумента через взвешенный семафор (зерно 1 МБ, вес файла по `os.Stat`). Самодостаточен: `go run solve.go <in> <out> [MB]`.
