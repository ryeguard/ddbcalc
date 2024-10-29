package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ryeguard/ddbcalc"
)

type Item struct {
	ID    string `dynamodbav:"id"`
	Name  string `dynamodbav:"name"`
	Age   int    `dynamodbav:"age"`
	Email string `dynamodbav:"email"`
}

const sizeLimit = 100

func main() {

	// A mock client is defined in the mock.go file.
	// The mock client is used to make this example run without needing to connect to a real DynamoDB table.
	client := newMockClient()

	// A real client could be created like this, and used as-is with the PutItem function below.
	// client := dynamodb.New(dynamodb.Options{Region: "us-west-2"})

	items := []Item{
		{
			ID:    "123",
			Name:  "John Doe",
			Age:   30,
			Email: "john@example.com",
		},
		{
			ID:    "456",
			Name:  "Pippilotta Delicatessa Windowshade Mackrelmint Efraimsdotter Longstocking",
			Age:   9,
			Email: "pippilotta.delicatessa.windowshade.mackrelmint.efraimdotter.longstocking@villavillekulla.se",
		},
	}

	for _, item := range items {
		m, err := attributevalue.MarshalMap(item)
		if err != nil {
			panic("failed to marshal item")
		}

		size, err := ddbcalc.MapSizeInBytes(m)
		if err != nil {
			panic("failed to calculate size")
		}

		fmt.Println("Item size:", size, "bytes")
		if size > sizeLimit {
			fmt.Printf("Item size is too large. Skipping item ID %v.", item.ID)
			continue
		} else {
			fmt.Printf("Item size is within limits. Writing item ID %v.", item.ID)
		}

		PutItem(context.TODO(), client, &dynamodb.PutItemInput{
			Item:      m,
			TableName: aws.String("my-table"),
		})
	}
}
