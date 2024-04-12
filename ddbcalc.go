package ddbcalc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ryeguard/ddbcalc/internal/calc"
)

const SizeLimitInBytes = 400_000 // 400 KB

// MapSizeInBytes returns the size of a map of AttributeValue in bytes.
func MapSizeInBytes(m map[string]types.AttributeValue) (int, error) {
	size := 0
	for k, v := range m {
		size += len([]byte(k))

		s, err := calc.SizeInBytes(&v)
		if err != nil {
			return 0, fmt.Errorf("size in bytes: %w", err)
		}

		size += s
	}

	return size, nil
}

// StructSizeInBytes returns the size of a struct in bytes.
// It is a convenience function as an alternative to first calling attributevalue.MarshalMap and then MapSizeInBytes.
func StructSizeInBytes(s interface{}) (int, error) {
	av, err := attributevalue.MarshalMap(s)
	if err != nil {
		return 0, fmt.Errorf("marshal map: %w", err)
	}

	return MapSizeInBytes(av)
}
