package main

import (
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
    count := 800
    reqcount := 3
    url := "http://127.0.0.1:8000"
    mode := "GET"

    var wg sync.WaitGroup

    for i := 0; i < count; i++ {
        go func() {
            defer wg.Done()
            Flood(reqcount, url, mode)
        }()
        wg.Add(1)
    }

    wg.Wait()
    fmt.Println("\nDone")
}
