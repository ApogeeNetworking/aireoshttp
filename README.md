# Cisco Wireless LAN Controller 
I do not advise using this as it uses the WLC frontend endpoints for authentication and data retrieval (but I'd rather use this than SSH)

This integration was testing on a WLC running version 8.5

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
    aps, err := cwlc.GetAps()
    if err != nil {
        log.Fatalf("%v", err)
    }
    fmt.Println(aps)
}
```
