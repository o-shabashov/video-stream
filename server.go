package main

import (
    "html"
    "io"
    "log"
    "net/http"
)

func main() {
	runServer()
}

func runServer() {
    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        log.Printf("%q %q", html.EscapeString(request.Method), html.EscapeString(request.URL.Path))
        _, _ = io.WriteString(writer, html.EscapeString(request.URL.Path))
    })

    log.Fatal(http.ListenAndServe(":8899", nil))
}