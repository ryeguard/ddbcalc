package ddbcalc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const SizeLimitInBytes = 400_000 // 400 KB

const (
	overheadMemberM = 3 // 3 byte
	overheadMemberL = 3 // 3 byte
	overheadElement = 1 // 1 byte
)

// MapSizeInBytes returns the size of a map of AttributeValue in bytes.
func MapSizeInBytes(m map[string]types.AttributeValue) (int, error) {
	size := 0
	for k, v := range m {
		size += len([]byte(k))
		s, err := sizeInBytes(&v)
		if err != nil {
			return 0, err
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
		return 0, err
	}

	return MapSizeInBytes(av)
}

func sizeInBytes(av *types.AttributeValue) (int, error) {
	if av == nil {
		return 0, nil
	}

	switch _av := (*av).(type) {
	case nil:
		return 0, nil
	case *types.AttributeValueMemberS:
		return len(_av.Value), nil
	case *types.AttributeValueMemberN:
		return len(_av.Value), nil
	case *types.AttributeValueMemberB:
		return len(_av.Value), nil
	case *types.AttributeValueMemberBOOL:
		return 1, nil
	case *types.AttributeValueMemberNULL:
		return 1, nil
	case *types.AttributeValueMemberBS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}
		return size, nil
	case *types.AttributeValueMemberL:
		size := overheadMemberL
		for _, v := range _av.Value {
			s, err := sizeInBytes(&v)
			if err != nil {
				return 0, err
			}
			size += s
			size += overheadElement
		}
		return size, nil
	case *types.AttributeValueMemberM:
		size := overheadMemberM
		for k, v := range _av.Value {
			size += len(k)
			s, err := sizeInBytes(&v)
			if err != nil {
				return 0, err
			}
			size += s
			size += overheadElement
		}
		return size, nil

	case *types.AttributeValueMemberNS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}
		return size, nil
	case *types.AttributeValueMemberSS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}
		return size, nil
	default:
		return 0, fmt.Errorf("unknown type: %T", _av)
	}
}
