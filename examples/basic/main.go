package main

import (
	"fmt"

	"github.com/ryeguard/ddbcalc"
)

type Item struct {
	ID    string
	Name  string
	Age   int
	Email string
}

func main() {
	item := Item{
		ID:    "123",
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	size, err := ddbcalc.StructSizeInBytes(item)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Item size: %d bytes\n", size)
}
