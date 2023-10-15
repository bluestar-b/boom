package main

import (
    "flag"
    "fmt"
    "sync"
    "github.com/valyala/fasthttp"

)


func CreateRequest(url string, method string) ([]byte, int) {
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(method)
    req.SetRequestURI(url)

    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    if err := fasthttp.Do(req, resp); err != nil {
        return []byte{}, 0
    }

    return resp.Body(), resp.StatusCode()
}


func Flood(reqcount int, url string, mode string) {
    for i := 0; i < reqcount; i++ {
        _, getStatusCode := CreateRequest(url, mode)
        fmt.Printf("status: %d, %s\n", getStatusCode, url)
    }

}

func main() {
    countPtr := flag.Int("count", 8, "Number of goroutines")
    reqCountPtr := flag.Int("reqcount", 3, "Number of requests per goroutine")
    urlPtr := flag.String("url", "", "URL to request")
    modePtr := flag.String("mode", "", "HTTP request method")
    flag.Parse()

    if *urlPtr == "" || *modePtr == "" {
        fmt.Println("Please provide both a URL and an HTTP request method.")
        return
    }

    var wg sync.WaitGroup

    for i := 0; i < *countPtr; i++ {
        go func() {
            defer wg.Done()
            Flood(*reqCountPtr, *urlPtr, *modePtr)
        }()
        wg.Add(1)
    }

    wg.Wait()
    fmt.Println("\nDone")
}

