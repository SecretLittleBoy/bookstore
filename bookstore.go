package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

const (
	defaultCursor   = "0"
	defaultPageSize = 2
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
		return nil, status.Error(codes.Internal, "failed to list shelves")
	}
	shelfListPb := make([]*pb.Shelf, 0, len(shelfListDatabase))
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
		return nil, status.Error(codes.Internal, "failed to create shelf")
	}
	return &pb.Shelf{Id: shelfDatabase.ID, Theme: shelfDatabase.Theme, Size: shelfDatabase.Size}, nil
}

func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	shelfDatabase, err := s.bookstore.GetShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get shelf")
	}
	return &pb.Shelf{Id: shelfDatabase.ID, Theme: shelfDatabase.Theme, Size: shelfDatabase.Size}, nil
}

func (s *server) DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	err := s.bookstore.DeleteShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete shelf")
	}
	return &emptypb.Empty{}, nil
}

func (s *server) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	var (
		cursor   string = defaultCursor
		pageSize int    = defaultPageSize
	)
	if len(in.GetPageToken()) > 0 {
		pageInfo := Token(in.GetPageToken()).Decode()
		if !pageInfo.IsValid() {
			return nil, status.Error(codes.InvalidArgument, "invalid page token")
		}
		cursor = pageInfo.NextID
		pageSize = int(pageInfo.PageSize)
	}
	bookList, err := s.bookstore.GetBookListByShelfID(ctx, in.GetShelf(), cursor, pageSize+1)
	if err != nil {
		fmt.Printf("GetBookListByShelfID failed, err:%v\n", err)
		return nil, status.Error(codes.Internal, "query failed")
	}
	var (
		hasNextPage   bool
		nextPageToken string
		realSize      int = len(bookList)
	)
	if len(bookList) > pageSize {
		hasNextPage = true
		realSize = pageSize
	}
	res := make([]*pb.Book, 0, len(bookList))
	for i := 0; i < realSize; i++ {
		res = append(res, &pb.Book{
			Id:     bookList[i].ID,
			Author: bookList[i].Author,
			Title:  bookList[i].Title,
		})
	}
	if hasNextPage {
		nextPageInfo := Page{
			NextID:        strconv.FormatInt(res[realSize-1].Id, 10), // res[realSize-1].Id 最后一个返回结果的id
			NextTimeAtUTC: time.Now().Unix(),
			PageSize:      int64(pageSize),
		}
		nextPageToken = string(nextPageInfo.Encode())
	}
	return &pb.ListBooksResponse{
		Books:         res,
		NextPageToken: nextPageToken,
	}, nil
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
