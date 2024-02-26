package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	useSamePort = false
)

func main() {
	db, err := NewDB("bookstore.db")
	if err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
		return
	}

	server := server{
		bookstore: &bookstore{db: db},
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("listen failed,err:%v\n", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterBookstoreServer(s, &server)

	if useSamePort {// 使用同一个端口
		gwmux := runtime.NewServeMux() // 1. 创建gRPC-Gateway mux
		dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		if err := pb.RegisterBookstoreHandlerFromEndpoint(
			context.Background(), gwmux, "127.0.0.1:8080", dops); err != nil {
			fmt.Printf("RegisterBookstoreHandlerFromEndpoint failed, err:%v\n", err)
			return
		}

		mux := http.NewServeMux() // 2. 新建HTTP mux
		mux.Handle("/", gwmux)

		gwServer := &http.Server{ // 3. 定义HTTP Server
			Addr:    "127.0.0.1:8080",
			Handler: grpcHandlerFunc(s, mux),
		}
		// 4. 启动服务
		fmt.Println("serving on 127.0.0.1:8080...")
		gwServer.Serve(l)
	} else {// 使用不同的端口
		go func() {
			if err := s.Serve(l); err != nil {
				fmt.Printf("failed to serve,err:%v\n", err)
				return
			}
		}()

		conn, err := grpc.DialContext(context.Background(), "localhost:8080",
			grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("grpc.DialContext failed,err:%v\n", err)
			return
		}
		defer conn.Close()

		gwmux := runtime.NewServeMux()
		err = pb.RegisterBookstoreHandler(context.Background(), gwmux, conn)
		if err != nil {
			fmt.Printf("pb.RegisterBookstoreHandler failed,err:%v\n", err)
			return
		}
		gwServer := &http.Server{
			Addr:    ":8081",
			Handler: gwmux,
		}
		fmt.Println("grpc-Gateway serve on :8081...")
		if err := gwServer.ListenAndServe(); err != nil {
			fmt.Printf("gwServer.ListenAndServe failed,err:%v\n", err)
			return
		}
	}
}

// grpcHandlerFunc 将gRPC请求和HTTP请求分别调用不同的handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
