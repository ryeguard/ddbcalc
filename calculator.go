package ddbcalc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func sizeInBytes(av *types.AttributeValue) (int, error) {
	if av == nil {
		return 0, nil
	}

	switch _av := (*av).(type) {
	case *types.AttributeValueMemberS:
		return len([]byte(_av.Value)), nil
	case *types.AttributeValueMemberN:
		return len([]byte(_av.Value)), nil
	case *types.AttributeValueMemberB:
		return len(_av.Value), nil
	case *types.AttributeValueMemberBOOL, *types.AttributeValueMemberNULL:
		return 1, nil
	case *types.AttributeValueMemberBS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}
		return size, nil
	case *types.AttributeValueMemberL:
		size := 3
		for _, v := range _av.Value {
			s, err := sizeInBytes(&v)
			if err != nil {
				return 0, err
			}
			size += s
			size++
		}
		return size, nil
	case *types.AttributeValueMemberM:
		size := 3
		for k, v := range _av.Value {
			size += len([]byte(k))
			s, err := sizeInBytes(&v)
			if err != nil {
				return 0, err
			}
			size += s
			size++
		}
		return size, nil

	case *types.AttributeValueMemberNS:
		size := 0
		for _, v := range _av.Value {
			size += len([]byte(v))
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

func StructSizeInBytes(s interface{}) (int, error) {
	av, err := attributevalue.MarshalMap(s)
	if err != nil {
		return 0, err
	}

	size := 0
	for k, v := range av {
		size += len([]byte(k))
		s, err := sizeInBytes(&v)
		if err != nil {
			return 0, err
		}
		size += s
	}
	return size, nil
}
