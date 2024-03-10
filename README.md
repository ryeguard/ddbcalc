# ddbcalc

[![Go Coverage](https://github.com/ryeguard/ddbcalc/wiki/coverage.svg)](https://raw.githack.com/wiki/ryeguard/ddbcalc/coverage.html)

A DynamoDB item size calculator.

This package has no dependencies other than AWS's [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2).

## Raison d'être

DynamoDB has a limit of 400KB for the size of an item. The `ddbcalc` package provides functionality to calculate the size of any Go struct in bytes. This can be useful when you need to ensure that the size of an item does not exceed the limit, either through tests or at runtime.

For more info, please see the blog post at [rygard.se](https://www.rygard.se/blog/240220_dynamodb_item_size/240220_dynamodb_item_size.html).

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
		fmt.Println(err)
		return
	}

	fmt.Printf("Item size: %d bytes\n", size)
}
```

You can run the example above with the following command:

```sh
go run example/basic/main.go
```

For more examples, see the [examples](./examples) directory and also check out the tests in all `*_test.go` files.

## Contributing

Contributions are welcome in any shape or form. Please feel free to open an issue or a pull request.
