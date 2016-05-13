package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/oleksandr/bonjour"
)

var macRE *regexp.Regexp

func init() {
	macRE = regexp.MustCompile("([A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2})")
}

func main() {
	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize resolver: %v\n", err)
		os.Exit(1)
	}

	results := make(chan *bonjour.ServiceEntry)

	go func(exitCh chan<- bool) {
		for e := range results {
			fmt.Printf("%s\n", e.Instance)
			fmt.Printf("RE findall: %s", macRE.FindString(e.Instance))
			exitCh <- true
			time.Sleep(1 * time.Second)
			//			os.Exit(0)
		}
	}(resolver.Exit)

	err = resolver.Browse("_workstation._tcp", "local.", results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to browse: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://api.macvendors.com/84:D6:D0:4F:EF:BA")

	if err != nil {
		fmt.Printf("Got an error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("We got the response: %s\n", resp)

	fmt.Printf("Our status is %s (code %d)\n", resp.Status, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Our body is %s\n", body)
}
