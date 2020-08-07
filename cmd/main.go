package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	aireoshttp "github.com/drkchiloll/cisco-wlc"
	"github.com/subosito/gotenv"
)

var host, user, pass string

func init() {
	gotenv.Load()
	host = os.Getenv("SSH_HOST")
	user = os.Getenv("SSH_USER")
	pass = os.Getenv("SSH_PW")
}

func main() {
	c := aireoshttp.New(host, user, pass, true)

	err := c.Login()
	if err != nil {
		log.Fatalf("%v", err)
	}
	aps, _ := c.GetAps()

	// Concurrency Example
	apDetail := make(chan aireoshttp.ApDetail)
	go proc(aps, c, apDetail)
	for ap := range apDetail {
		fmt.Println(ap)
	}
	// for _, ap := range aps {
	// 	d, _ := c.GetApDetails(ap.MacAddr)
	// 	fmt.Println(d.IPAddr)
	// }
}

func proc(aps []aireoshttp.AP, c *aireoshttp.Client, a chan aireoshttp.ApDetail) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	for _, ap := range aps {
		wg.Add(1)
		sem <- struct{}{}
		go func(ap aireoshttp.AP) {
			defer func() {
				<-sem
				wg.Done()
			}()
			d, err := c.GetApDetails(ap.MacAddr)
			if err != nil {
				fmt.Println(err)
			}
			a <- d
		}(ap)
	}
	wg.Wait()
	close(a)
}
