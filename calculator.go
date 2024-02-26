package ddbcalc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Do we need to export this see, comment below.
func SizeInBytes(av *types.AttributeValue) (int, error) {
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
		// Why is the size 3 here? THis might be worth a comment.
		size := 3
		for _, v := range _av.Value {
			s, err := SizeInBytes(&v)
			if err != nil {
				return 0, err
			}
			size += s
			// Why do we plus 1 here?
			size++
		}
		return size, nil
	case *types.AttributeValueMemberM:
		// Same comment as above
		size := 3
		for k, v := range _av.Value {
			size += len([]byte(k))
			s, err := SizeInBytes(&v)
			if err != nil {
				return 0, err
			}
			size += s
			// Why do we plus 1 here?
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

// Is size enough as name?
// I think name could be shorter if you add some a comment with some docs.
func StructSizeInBytes(s interface{}) (int, error) {
	av, err := attributevalue.MarshalMap(s)
	if err != nil {
		return 0, err
	}

	return MapSizeInBytes(av)
}

// Do you expect a user to access this function if not, do you want to export it?
// Less exporter MIGHT make your life easier since you have not given access for users to in
// so you can change it!
func MapSizeInBytes(m map[string]types.AttributeValue) (int, error) {
	size := 0
	for k, v := range m {
		size += len([]byte(k))
		s, err := SizeInBytes(&v)
		if err != nil {
			return 0, err
		}
		size += s
	}
	return size, nil
}
