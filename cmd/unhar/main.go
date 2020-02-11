package main

import (
	"fmt"
	"os"

	"arp242.net/har"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "unhar: must give at least one filename")
		os.Exit(1)
	}

	for _, f := range os.Args[1:] {
		h, err := har.FromFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unhar: reading %q: %s\n", f, err)
			os.Exit(1)
		}

		err = h.Extract(false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unhar: extracting %q: %s\n", f, err)
			os.Exit(1)
		}
	}
}
