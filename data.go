package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
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
