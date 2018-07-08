package client

import (
        "context"
	"fmt"
        "google.golang.org/grpc"
        "github.com/tortuoise/cacheservice/rpc"
        "os"
        "time"
)


func ClientMain() {
    if err := RunClient(context.Background()); err != nil {
        fmt.Fprintf(os.Stderr, "failed: %v\n", err)
        os.Exit(1)
    }
}

// RunClient Should just return the "raw" error from rpc calls
func RunClient(ctx context.Context) error {
    // connnect
    conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(ClientInterceptor))
    if err != nil {
        return fmt.Errorf("failed to dial server: %v", err)
    }
    cache := rpc.NewCacheClient(conn)
    // store
    _, err = cache.Store(ctx, &rpc.StoreReq{AccountToken: "doofus", Key: "gopher", Val: []byte("con")})
    if err != nil {
        return err
    }
    // get
    resp, err := cache.Get(ctx, &rpc.GetReq{Key: "gopher"})
    if err != nil {
        return err
    }
    fmt.Printf("Got cached value %s\n", resp.Val)
    // get NotFound
    resp, err = cache.Get(ctx, &rpc.GetReq{Key: "opher"})
    if err != nil {
        return err
    }
    fmt.Printf("Got cached value %s\n", resp.Val)
    return nil
}

func ClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
        start := time.Now()
        err := invoker(ctx, method, req, reply, cc, opts...)
        fmt.Fprintf(os.Stderr, "invoke remote method=%s duration=%v error=%v", method, time.Since(start)/1e6, err)
        return err
}
