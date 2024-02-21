package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"net"
	"net/http"

	 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	go func() {
		if err := s.Serve(l); err != nil {
			fmt.Printf("failed to serve,err:%v\n", err)
			return
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "localhost:8080",
		grpc.WithBlock(),grpc.WithTransportCredentials(insecure.NewCredentials()))
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
