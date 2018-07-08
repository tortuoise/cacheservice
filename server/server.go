package server

import (
        "context"
	"fmt"
        "github.com/tortuoise/cacheservice/rpc"
        "google.golang.org/grpc"
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/status"
        "math/rand"
        "net"
        "os"
        "time"
)

func init() {

}

func ServerMain() {
    if err := RunServer(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to run cache server: %s\n", err)
        os.Exit(1)
    }
}

func RunServer() error {
        srv := grpc.NewServer()
        ms := &CacheService{store: make(map[string][]byte), keysByAccount: make(map[string]int64)}
        rpc.RegisterCacheServer(srv, ms)
        rpc.RegisterAccountsServer(srv, ms)
        l, err := net.Listen("tcp", "localhost:5051")
        if err != nil {
                return err
        }
        conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
        if err != nil {
                return err
        }
        ms.accounts = rpc.NewAccountsClient(conn)
        // blocks until complete
        return srv.Serve(l)
}

type CacheService struct {
        accounts rpc.AccountsClient
        store map[string][]byte
        keysByAccount map[string]int64
}

func (s *CacheService) Get(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
        abortChan := make(chan bool, 1)
        errChan := make(chan error, 1)
        respChan := make(chan []byte, 1)

        go func() {
                <-ctx.Done()
                abortChan<-true
        }()
        go func() {
                val, ok := s.store[req.Key]
                if !ok {
                        errChan<- status.Errorf(codes.NotFound, "Key not found %s",
                    req.Key)
                } else {
                        respChan<- val
                }
        }()
        time.Sleep(time.Duration(rand.Intn(2000))*time.Millisecond)
        for {
                select {
                case <-abortChan:
                        fmt.Fprintf(os.Stderr, "Canceled: %v\n", req.Key)
                        return nil, status.Errorf(codes.Canceled, "Canceled")
                case err := <-errChan:
                        fmt.Fprintf(os.Stderr, "NotFound: %v\n", req.Key)
                        return nil, err
                case val := <-respChan:
                        fmt.Fprintf(os.Stderr, "Responding: %v:%v\n", req.Key, val)
                        return &rpc.GetResp{Val: val}, nil
                }
        }
}

func (s *CacheService) Store(ctx context.Context, req *rpc.StoreReq) (*rpc.StoreResp,
    error) {
        if s.store == nil || s.keysByAccount == nil {
               return nil, status.Errorf(codes.FailedPrecondition, "Store not set up %v\n", req.AccountToken)
        }
        resp, err := s.accounts.GetByToken(ctx, &rpc.GetByTokenReq{req.AccountToken})
        if err != nil {
                return nil, err
        }
        if s.keysByAccount[req.AccountToken] > resp.Account.MaxCacheKeys {
                return nil, status.Errorf(codes.FailedPrecondition, "Account %s exceeds max key limit %d", req.AccountToken, resp.Account.MaxCacheKeys)
        }
        if _, ok := s.store[req.Key]; !ok {
                s.keysByAccount[req.AccountToken] += 1
        }
        s.store[req.Key] = req.Val
        return &rpc.StoreResp{}, nil
        //return nil, fmt.Errorf("unimplemented")
}

func (s *CacheService) GetByToken(ctx context.Context, req *rpc.GetByTokenReq) (*rpc.GetByTokenResp, error) {
        //if _, ok := s.keysByAccount[req.Token]; !ok {
                return &rpc.GetByTokenResp{&rpc.Account{100}}, nil
        //}
        //return &rpc.GetByTokenResp{}, status.Errorf(codes.NotFound, "Token not found
}
