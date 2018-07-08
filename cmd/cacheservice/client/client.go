package main

import (
        "context"
	"fmt"
        "github.com/tortuoise/cacheservice/client"
        "google.golang.org/grpc/status"
        "google.golang.org/grpc/codes"
        "os"
        "os/signal"
        "syscall"
        "time"
)

func main() {

        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
        go func(chan os.Signal) {
                <-ch
                cancel()
        }(ch)
        if err := client.RunClient(ctx); err != nil {
                st := status.Code(err)
                switch st {
                case codes.Canceled:
                        fmt.Fprintf(os.Stderr, "Canceled: %v\n", err)
                        os.Exit(1)
                case codes.DeadlineExceeded:
                        fmt.Fprintf(os.Stderr, "DeadlineExceeded: %v\n", err)
                        os.Exit(1)
                case codes.NotFound:
                        fmt.Fprintf(os.Stderr, "NotFound: %v\n", err)
                        os.Exit(1)
                default:
                        fmt.Fprintf(os.Stderr, "Failed: %v %v\n", st, err)
                        os.Exit(1)
                }
        }

}
