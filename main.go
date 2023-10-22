package main

import (
    "flag"
    "fmt"
    "sync"
    "time"
    "github.com/valyala/fasthttp"
)

func CreateRequest(url string, method string) (int, time.Duration) {
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    req.Header.SetMethod(method)
    req.SetRequestURI(url)

    start := time.Now()
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    if err := fasthttp.Do(req, resp); err != nil {
        return 0, time.Since(start)
    }

    return resp.StatusCode(), time.Since(start)
}

func Flood(reqcount int, url string, mode string) (int, int) {
    successCount := 0
    failureCount := 0

    for i := 0; i < reqcount; i++ {
        getStatusCode, requestTime := CreateRequest(url, mode)
        if getStatusCode == 200 {
            successCount++
        } else {
            failureCount++
        }
        fmt.Printf("Request %d - Status: %d, Time: %s\n", i+1, getStatusCode, requestTime)
    }

    return successCount, failureCount
}

func main() {
    countPtr := flag.Int("c", 8, "Number of goroutines")
    reqCountPtr := flag.Int("rc", 1, "Number of requests per goroutine")
    urlPtr := flag.String("u", "", "URL to request")
    modePtr := flag.String("m", "", "HTTP request method")
    flag.Parse()

    if *urlPtr == "" || *modePtr == "" {
        fmt.Println("Please provide both a URL and an HTTP request method.")
        return
    }

    var wg sync.WaitGroup
    totalSuccess := 0
    totalFailure := 0

    for i := 0; i < *countPtr; i++ {
        go func() {
            defer wg.Done()
            success, failure := Flood(*reqCountPtr, *urlPtr, *modePtr)
            totalSuccess += success
            totalFailure += failure
        }()
        wg.Add(1)
    }

    wg.Wait()

    totalRequests := *countPtr * *reqCountPtr
    fmt.Printf("\nTotal Requests: %d\n", totalRequests)
    fmt.Printf("Successful Requests: %d (%.2f%%)\n", totalSuccess, float64(totalSuccess)/float64(totalRequests)*100)
    fmt.Printf("Failed Requests: %d (%.2f%%)\n", totalFailure, float64(totalFailure)/float64(totalRequests)*100)
}

