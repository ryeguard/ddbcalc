package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBPutItemAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func PutItem(ctx context.Context, api DynamoDBPutItemAPI, input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return api.PutItem(ctx, input)
}

// mockPutItemClient is a mock implementation of the DynamoDBPutItemAPI interface.
type mockPutItemClient struct {
	mockPutItemAPI
}

// newMockClient returns a new mock client.
// The mock client may be used interchangeably with the real client for PutItem operations.
func newMockClient() *mockPutItemClient {
	return &mockPutItemClient{mockPutItemAPI: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		fmt.Println("Mocked PutItem called.")
		if params.TableName == nil {
			return nil, fmt.Errorf("TableName is nil")
		}
		fmt.Println("Writing to table:", *params.TableName)
		if params.Item == nil {
			return nil, fmt.Errorf("Item is nil")
		}

		for k := range params.Item {
			fmt.Printf("\tItem key: %v\n", k)
		}
		return &dynamodb.PutItemOutput{}, nil
	}}
}

// mockPutItemAPI is a type that represents a function that can be mocked.
type mockPutItemAPI func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)

// PutItem is a mock implementation of the DynamoDBPutItemAPI interface.
func (m mockPutItemAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m(ctx, params, optFns...)
}
