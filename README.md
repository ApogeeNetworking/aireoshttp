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
    // Get AP Details Concurrently
    var wg sync.WaitGroup
    // Throttle the Number of Simultaneous Requests
    sem := make(chan struct{}, 5)
    var apDetails []cwlc.ApDetail
    for _, ap := range aps {
        sem <- struct{}{} // Add to the Semaphore
        wg.Add(1)
        go func(ap cwlc.AP) {
            defer func() {
                <-sem // Release the Resource
                wg.Done()
            }()
            d, _ := wlc.GetApDetails(ap.MacAddr)
            apDetails = append(apDetails, d)
        }(ap)
    }
    wg.Wait()
    /**
        ApDetail = Name, MacAddr, IPAddr, RemoteSW, RemoteIntf ...
        Refer to cwlc.go for Type Declaration 
    */
}
```
