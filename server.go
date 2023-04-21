package main

import (
    "fmt"
    "html"
    "log"
    "net/http"
    "strconv"
    "sync"
)

var counter int
var rwMutex = &sync.RWMutex{}

func echoString(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
    rwMutex.Lock()
    counter++
    fmt.Fprintf(w, strconv.Itoa(counter))
    rwMutex.Unlock()
}

func readCounter(w http.ResponseWriter, r *http.Request) {
    rwMutex.RLock()
    defer rwMutex.RUnlock()
    fmt.Fprintf(w, strconv.Itoa(counter))
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", echoString)
    mux.HandleFunc("/increment", incrementCounter)
    mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi")
    })
    mux.HandleFunc("/counter", readCounter)

    log.Fatal(http.ListenAndServe(":8081", mux))
}
