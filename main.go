package main

import (
	"fmt"
)

func main() {
	db, err := NewDB("bookstore.db")
	if err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
		return
	}
	



}
