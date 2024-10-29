package calc

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	overheadMemberM = 3 // 3 byte
	overheadMemberL = 3 // 3 byte
	overheadElement = 1 // 1 byte
)

func SizeInBytes(av *types.AttributeValue) int {
	if av == nil {
		return 0
	}

	switch _av := (*av).(type) {
	case *types.AttributeValueMemberS:
		return len(_av.Value)
	case *types.AttributeValueMemberN:
		return len(_av.Value)
	case *types.AttributeValueMemberB:
		return len(_av.Value)
	case *types.AttributeValueMemberBOOL:
		return 1
	case *types.AttributeValueMemberNULL:
		return 1
	case *types.AttributeValueMemberBS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}

		return size
	case *types.AttributeValueMemberL:
		return listSize(_av)
	case *types.AttributeValueMemberM:
		return mapSize(_av)
	case *types.AttributeValueMemberNS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}

		return size
	case *types.AttributeValueMemberSS:
		size := 0
		for _, v := range _av.Value {
			size += len(v)
		}

		return size
	default:
		return 0
	}
}

func listSize(l *types.AttributeValueMemberL) int {
	size := overheadMemberL

	for _, v := range l.Value {
		size += SizeInBytes(&v)
		size += overheadElement
	}

	return size
}

func mapSize(m *types.AttributeValueMemberM) int {
	size := overheadMemberM

	for k, v := range m.Value {
		size += len(k)

		size += SizeInBytes(&v)
		size += overheadElement
	}

	return size
}
