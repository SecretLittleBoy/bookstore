package main

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultShelfSize = 10
)

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Shelf{}, &Book{})
	return db, nil
}

type Shelf struct {
	ID       int64 `gorm:"primaryKey"`
	Theme    string
	Size     int64
	CreatAt  time.Time
	UpdateAt time.Time
}

type Book struct {
	ID       int64 `gorm:"primaryKey"`
	Author   string
	Title    string
	ShelfID  int64
	CreatAt  time.Time
	UpdateAt time.Time
}

type bookstore struct {
	db *gorm.DB
}

func (b *bookstore) CreateShelf(ctx context.Context, shelf Shelf) (*Shelf, error) {
	if len(shelf.Theme) <= 0 {
		return nil, errors.New("invalid theme")
	}
	if shelf.Size <= 0 {
		shelf.Size = defaultShelfSize
	}
	v := Shelf{Theme: shelf.Theme, Size: shelf.Size, CreatAt: time.Now(), UpdateAt: time.Now()}
	if err := b.db.Create(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (b *bookstore) GetShelf(ctx context.Context, id int64) (*Shelf, error) {
	var v Shelf
	if err := b.db.WithContext(ctx).First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (b *bookstore) ListShelf(ctx context.Context) ([]*Shelf, error) {
	var vs []*Shelf
	if err := b.db.WithContext(ctx).Find(&vs).Error; err != nil {
		return nil, err
	}
	return vs, nil
}

func (b *bookstore) DeleteShelf(ctx context.Context, id int64) error {
	return b.db.WithContext(ctx).Delete(&Shelf{}, id).Error
}

func (b *bookstore) GetBookListByShelfID(ctx context.Context, shelfID int64, cursor string, pageSize int) ([]*Book, error) {
	var vl []*Book
	err := b.db.WithContext(ctx).Where("shelf_id = ? and id > ?", shelfID, cursor).Order("id asc").Limit(pageSize).Find(&vl).Error
	return vl, err
}

