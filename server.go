package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	sampleVideo    = "sample.mp4"
	fileBufferSize = 1024
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%q %q", request.Method, request.URL.Path)

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

func copyFile(writer http.ResponseWriter, file *os.File) (written int64) {
	reader := bufio.NewReader(file)
	fileBuffer := make([]byte, fileBufferSize)
	written, err := io.CopyBuffer(writer, reader, fileBuffer)
	check(err)

	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
