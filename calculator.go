package ddbcalc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func sizeInBytes(av *types.AttributeValue) int {
	if av == nil {
		return 0
	}

	switch _av := (*av).(type) {
	case *types.AttributeValueMemberS:
		return len([]byte(_av.Value))
	case *types.AttributeValueMemberN:
		return len([]byte(_av.Value))
	case *types.AttributeValueMemberB:
		return len(_av.Value)
	case *types.AttributeValueMemberBOOL, *types.AttributeValueMemberNULL:
		return 1
	case *types.AttributeValueMemberBS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}
		return size
	case *types.AttributeValueMemberL:
		size := 3
		for _, v := range _av.Value {
			size += sizeInBytes(&v)
			size++
		}
		return size
	case *types.AttributeValueMemberM:
		size := 3
		for k, v := range _av.Value {
			size += len([]byte(k))
			size += sizeInBytes(&v)
			size++
		}
		return size

	case *types.AttributeValueMemberNS:
		size := 0
		for _, v := range _av.Value {
			size += len([]byte(v))
		}
		return size
	case *types.AttributeValueMemberSS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}
		return size
	default:
		panic(fmt.Sprintf("unknown type %T", _av))
	}
}

func StructSizeInBytes(s interface{}) int {
	av, err := attributevalue.MarshalMap(s)
	if err != nil {
		panic(err)
	}

	size := 0
	for k, v := range av {
		size += len([]byte(k))
		size += sizeInBytes(&v)
	}
	return size
}
