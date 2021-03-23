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
	sampleVideo       = "sample.mp4"
	chunkSize   int64 = 100 * 1024
)

var (
	part  []byte
	count int
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%q %q", html.EscapeString(request.Method), html.EscapeString(request.URL.Path))

		file, err := os.Open(sampleVideo)
		check(err)
		defer file.Close()

		fileHeader := make([]byte, 512)
		file.Read(fileHeader)
		contentType := http.DetectContentType(fileHeader)
		file.Seek(0, 0)

		movieFileStat, err := file.Stat()
		check(err)

		writer.Header().Set("Content-Type", contentType)
		writer.Header().Set("Content-Length", strconv.Itoa(int(movieFileStat.Size())))

		//channel := make(chan int64)

		reader := bufio.NewReader(file)
		//buffer := bytes.NewBuffer(make([]byte, 0))
		part = make([]byte, chunkSize)

		for {
			start := time.Now()
			written, err := io.CopyN(writer, reader, chunkSize)
			if err != nil {
				break
			}
			duration := time.Since(start)
			log.Printf("Written: %v", written)
			log.Printf("Duration: %v", duration)
		}

		if err != io.EOF && err != nil {
			log.Fatal("Error Reading ", sampleVideo, ": ", err)
		} else {
			err = nil
		}

		//for i := range channel {
		//	log.Println(i)
		//}

		return
	})

	log.Fatal(http.ListenAndServe(":8899", nil))
}

func copyFile(writer http.ResponseWriter, file *os.File) {
	fileStat, err := file.Stat()
	check(err)
	var i int64

	for i = 0; i < fileStat.Size(); i += 100 * 1024 {
		start := time.Now()
		_, err := io.CopyN(writer, file, 100*1024)
		check(err)
		file.Seek(i, 1)
		duration := time.Since(start)

		log.Println(duration)

		//channel <- duration.Milliseconds()
	}

	//close(channel)
}

func check(e error) {
	if e == io.EOF {
		return
	}
	if e != nil {
		panic(e)
	}
}
