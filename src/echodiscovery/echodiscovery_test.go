package echodiscovery

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestHasAmazonDevice(t *testing.T) {
	has, err := NetworkHasAmazonDevice(1 * time.Second)
	if err != nil {
		t.Fatalf("Not expecting error detecting presence of Amazon device: %v", err)
	}
	fmt.Fprintf(os.Stdout, "Amazon device present? %s", has)
}
