package calc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	overheadMemberM = 3 // 3 byte
	overheadMemberL = 3 // 3 byte
	overheadElement = 1 // 1 byte
)

func SizeInBytes(av *types.AttributeValue) (int, error) {
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
		return listSize(_av)
	case *types.AttributeValueMemberM:
		return mapSize(_av)
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

func listSize(l *types.AttributeValueMemberL) (int, error) {

	size := overheadMemberL
	for _, v := range l.Value {
		s, err := SizeInBytes(&v)
		if err != nil {
			return 0, err
		}

		size += s
		size += overheadElement
	}

	return size, nil
}

func mapSize(m *types.AttributeValueMemberM) (int, error) {

	size := overheadMemberM
	for k, v := range m.Value {
		size += len(k)

		s, err := SizeInBytes(&v)
		if err != nil {
			return 0, err
		}

		size += s
		size += overheadElement
	}

	return size, nil
}

func unused() {
	// This is a dummy function to test linting
}
