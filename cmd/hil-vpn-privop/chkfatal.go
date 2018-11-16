package main

import (
	"fmt"
	"os"
)

func chkfatal(ctx string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s: %v\n", ctx, err)
		os.Exit(1)
	}
}
