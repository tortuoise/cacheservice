package main

import (
	"fmt"
        "github.com/tortuoise/cacheservice/server"
        "os"
)

func main() {
        if err := server.RunServer(); err != nil {
                fmt.Fprintf(os.Stderr, "Failed to run cache server: %s\n", err)
                os.Exit(1)
        }
}
