/*
You can compile this on your Mac, then, copy to Linux

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
go build -o url url.go


*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func Get(url string) error {

	goodUrl := fmt.Sprintf("http://%s", url)

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5000 * time.Millisecond,
		}).Dial,
		TLSHandshakeTimeout: 5000 * time.Millisecond,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * 3,
		Transport: netTransport,
	}

	resp, err := netClient.Get(goodUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	if err != nil {
		fmt.Println(err)
	}
	return err

}

func main() {

	file := os.Args
	if len(file) <= 1 {
		log.Fatalf("Need to enter filename. Try the following:\n\n./url file\n\n")
	}
	fmt.Println(file, len(file))

	dat, err := ioutil.ReadFile(file[1])
	if err != nil {
		log.Fatalf("That wasn't a valid file")
	}
	records := strings.Split(string(dat), "\n")

	var wg sync.WaitGroup
	for i, record := range records {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			err := Get(url)
			if err != nil {
				fmt.Println(err)
			}

		}(record)
		fmt.Println(i, record)

		if i%300 == 0 && i != 0 {
			time.Sleep(10 * time.Second)
			log.Printf("\n\nWe need to sleep for 10 seconds\n\n")
		}
	}
	wg.Wait()

}
