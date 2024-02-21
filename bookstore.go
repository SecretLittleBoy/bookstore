package main

import (
	"bookstore/pb"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedBookstoreServer
	bookstore *bookstore
}

func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	shelfListDatabase, err := s.bookstore.ListShelf(ctx)
	if err != nil {
		if err == gorm.ErrEmptySlice {
			return &pb.ListShelvesResponse{}, nil
		}
		return nil, status.Error(codes.Internal, "failed to list shelves: %v")
	}
	shelfListPb := make([]*pb.Shelf, len(shelfListDatabase))
	for _, shelf := range shelfListDatabase {
		shelfListPb = append(shelfListPb, &pb.Shelf{
			Id:    shelf.ID,
			Theme: shelf.Theme,
			Size:  shelf.Size,
		})
	}
	return &pb.ListShelvesResponse{Shelves: shelfListPb}, nil
}

func (s *server) CreateShelf(ctx context.Context, in *pb.CreateShelfRequest) (*pb.Shelf, error) {
	if len(in.GetShelf().GetTheme()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "shelf theme is required")
	}
	data := Shelf{
		Theme: in.GetShelf().GetTheme(),
		Size:  in.GetShelf().GetSize(),
	}
	shelfDatabase, err := s.bookstore.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create shelf: %v")
	}
	return &pb.Shelf{Id: shelfDatabase.ID, Theme: shelfDatabase.Theme, Size: shelfDatabase.Size}, nil
}

func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	shelfDatabase, err := s.bookstore.GetShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get shelf: %v")
	}
	return &pb.Shelf{Id: shelfDatabase.ID, Theme: shelfDatabase.Theme, Size: shelfDatabase.Size}, nil
}

func (s *server) DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	err := s.bookstore.DeleteShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete shelf: %v")
	}
	return &emptypb.Empty{}, nil
}

func (s *server) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBooks not implemented")
}

func (s *server) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.Book, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBook not implemented")
}

func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.Book, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBook not implemented")
}

func (s *server) DeleteBook(ctx context.Context, in *pb.GetBookRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBook not implemented")
}
