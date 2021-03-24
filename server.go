package main

import (
	"bufio"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	sampleVideo         = "sample.mp4"
	chunkSize1Mb  int64 = 1024 * 1024
	chunkSizeMbit       = float64(chunkSize1Mb * 8)
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%q %q", html.EscapeString(request.Method), html.EscapeString(request.URL.Path))

		file, err := os.Open(sampleVideo)
		check(err)
		defer file.Close()

		contentType, contentLength := fileInfo(file, err)
		writer.Header().Set("Content-Type", contentType)
		writer.Header().Set("Content-Length", contentLength)

		copyFile(writer, file)
	})

	log.Fatal(http.ListenAndServe(":8899", nil))
}

func fileInfo(file *os.File, err error) (string, string) {
	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	check(err)

	contentType := http.DetectContentType(fileHeader)
	_, err = file.Seek(0, 0)
	check(err)

	movieFileStat, err := file.Stat()
	check(err)

	return contentType, strconv.Itoa(int(movieFileStat.Size()))
}

func copyFile(writer http.ResponseWriter, file *os.File) {
	reader := bufio.NewReader(file)

	for {
		start := time.Now()
		_, err := io.CopyN(writer, reader, chunkSize1Mb)
		duration := time.Since(start)
		log.Printf("%.2f seconds, %.2f Mb/s", duration.Seconds(), chunkSizeMbit/duration.Seconds())

		if err == io.EOF {
			break
		} else if err != io.EOF && err != nil {
			panic(err)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
