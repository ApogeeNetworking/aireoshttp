# Cisco Wireless LAN Controller 
I do not advise using this as it uses the WLC frontend endpoints for authentication and data retrieval (but I'd rather use this than SSH)

This integration was tested on a WLC running version 8.5

## Installation

Install via **go get**:

```shell
go get -u github.com/drkchiloll/cisco-wlc
```

## Usage
Basic usage can be found below

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    cwlc "github.com/drkchiloll/cisco-wlc"
)

func main() {
    // Used in Development for SelfSigned Certs
    ignoreSSL := true
    wlc := cwlc.New("host/ip", "user", "pass", ignoreSSL)

    err := wlc.Login()
    if err != nil {
        log.Fatalf("%v", err)
    }
    aps, err := cwlc.GetAps() // Returns cwlc.[]AP (MacAddr, Name)
    if err != nil {
        log.Fatalf("%v", err)
    }
    // Synchronous Example
    for _, ap := range aps {
        d, _ := c.GetApDetails(ap.MacAddr)
        fmt.Println("Synchronous call: " + d.IPAddr)
    }

    // Concurrency Example
    apDetail := make(chan cwlc.ApDetail)
    go proc(aps, c, apDetail)
    // This Blocks (sync) until all APDetails
    // Are Received and the Channel is Closed (!ok)
    for {
        apd, ok := <- apDetail
        if !ok {
            // Channel is Closed
            break
        }
        fmt.Println("Async: " + apd.IPAddr)
    }
    // Besides using an Infinite for loop like above
    // You can "range" over the Streaming Data from
    // The channel (which is more succinct)
    // Comment out the below if using the Above
    for ap := range apDetail {
        fmt.Println("Async Streaming: " + ap.IPAddr)
    }

    /**
        ApDetail = Name, MacAddr, IPAddr, RemoteSW, RemoteIntf ...
        Refer to cwlc.go for Type Declaration 
    */
}

func proc(aps []cwlc.AP, c *cwlc.Client, a chan cwlc.ApDetail) {
    var wg sync.WaitGroup
    // Limit Resource Usage by using the Semaphore Pattern
    // Here we are only allowing 8 HTTP REQs at a time
    sem := make(chan struct{}, 8)
    for _, aps := range aps {
        wg.Add(1)
        sem <- struct{}{}
        go func(ap cwlc.AP) {
            defer func() {
                // Release the Resource
                <-sem
                // Decrement our WaitGroup
                wg.Done()
            }()
            // Make HTTP Request
            apdetail, _ := c.GetApDetails(ap.MacAddr)
            // Receive APDetail onto our Channel
            // Which is a Pointer(*) to an APDetail
            a <- apdetail
        }(ap)
    }
    wg.Wait()
    // Close the Channel
    close(a)
}
```
