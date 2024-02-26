# ddbcalc

A DynamoDB item size calculator

## Usage

The `StructSizeInBytes` function calculates the size of a struct in bytes. It may be used as follows:

```go
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
```

You can run the example above with the following command:

```sh
go run example/basic/main.go
```

For more examples, see the [examples](./examples) directory and also check out the tests in all `*_test.go` files.
