package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%q %q", html.EscapeString(request.Method), html.EscapeString(request.URL.Path))

		movieFile, err := os.Open("D:\\Видео\\Game Records\\Clips\\GWF #1.mp4")
		if err != nil {
			panic(err)
		}
		defer movieFile.Close()

		fileHeader := make([]byte, 512)
		movieFile.Read(fileHeader)
		contentType := http.DetectContentType(fileHeader)

		movieFileStat, err := movieFile.Stat()
		if err != nil {
			panic(err)
		}

		writer.Header().Set("Content-Type", contentType)
		writer.Header().Set("Content-Length", strconv.Itoa(int(movieFileStat.Size())))

		movieFile.Seek(0, 0)

		_, err = io.Copy(writer, movieFile)
		if err != nil {
			panic(err)
		}
		return
	})

	log.Fatal(http.ListenAndServe(":8899", nil))
}
