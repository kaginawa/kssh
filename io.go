package main

import (
	"fmt"
	"io"
	"os"
)

func fatalf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, format+"\n", a); err != nil {
		fmt.Printf("DEBUG: "+format+"\n", a)
	}
	os.Exit(1)
}

func safeClose(closer io.Closer, name string) {
	if err := closer.Close(); err != nil {
		if err.Error() == "EOF" {
			return
		}
		if _, err := fmt.Fprintf(os.Stderr, "failed to close %s: %v\n", err, name); err != nil {
			fmt.Printf("failed to close %s: %v\n", err, name)
		}
	}
}

func handleError(_ int, err error) {
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "%v\n", err); err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}
