package main

import (
	"fmt"
	"os"
	"time"

	"github.com/oleksandr/bonjour"
)

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
			exitCh <- true
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}
	}(resolver.Exit)

	err = resolver.Browse("_workstation._tcp", "local.", results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to browse: %v\n", err)
		os.Exit(1)
	}

	select {}
}
