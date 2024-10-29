package calc

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestSizeInBytesNil(t *testing.T) {
	actual := SizeInBytes(nil)

	if actual != 0 {
		t.Errorf("got %d; want 0", actual)
	}
}

func Test_listSize(t *testing.T) {
	var tests = []struct {
		item     *types.AttributeValueMemberL
		name     string
		expected int
	}{
		{
			name:     "empty list",
			item:     &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
			expected: 3,
		},
		{
			name: "list",
			item: &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "abc"},
				&types.AttributeValueMemberN{Value: "123"},
				&types.AttributeValueMemberB{Value: []byte{1, 2, 3}},
			}},
			expected: 3 + // list overhead
				3*(3+1), // 3 elements 3 bytes each + 1 byte overhead
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := listSize(tc.item)

			if actual != tc.expected {
				t.Errorf("got %d; want %d", actual, tc.expected)
			}
		})
	}
}

func Test_mapSize(t *testing.T) {
	var tests = []struct {
		item     *types.AttributeValueMemberM
		name     string
		expected int
	}{
		{
			name:     "empty map",
			item:     &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{}},
			expected: 3,
		},
		{
			name: "map",
			item: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"s": &types.AttributeValueMemberS{Value: "abc"},
				"n": &types.AttributeValueMemberN{Value: "123"},
				"b": &types.AttributeValueMemberB{Value: []byte{1, 2, 3}},
			}},
			expected: 3 + // map overhead
				3*(1+3+1), // 3 elements 1 byte key + 3 bytes value + 1 byte overhead
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := mapSize(tc.item)

			if actual != tc.expected {
				t.Errorf("got %d; want %d", actual, tc.expected)
			}
		})
	}
}

func TestSizeInBytesOfBasicTypes(t *testing.T) {
	var tests = []struct {
		item     types.AttributeValue
		name     string
		expected int
	}{
		{
			name:     "nil",
			item:     nil,
			expected: 0,
		},
		{
			name:     "string",
			item:     &types.AttributeValueMemberS{Value: "abc"},
			expected: 3,
		},
		{
			name:     "number",
			item:     &types.AttributeValueMemberN{Value: "123"},
			expected: 3,
		},
		{
			name:     "binary",
			item:     &types.AttributeValueMemberB{Value: []byte{1, 2, 3}},
			expected: 3,
		},
		{
			name:     "bool",
			item:     &types.AttributeValueMemberBOOL{Value: true},
			expected: 1,
		},
		{
			name:     "null",
			item:     &types.AttributeValueMemberNULL{Value: true},
			expected: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SizeInBytes(&tc.item)

			if actual != tc.expected {
				t.Errorf("got %d; want %d", actual, tc.expected)
			}
		})
	}
}

func TestSizeInBytesOfSet(t *testing.T) {
	var tests = []struct {
		item     types.AttributeValue
		name     string
		expected int
	}{
		{
			name: "binary set",
			item: &types.AttributeValueMemberBS{Value: [][]byte{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}},
			expected: 9,
		},
		{
			name:     "number set",
			item:     &types.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
			expected: 3,
		},
		{
			name:     "string set",
			item:     &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}},
			expected: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SizeInBytes(&tc.item)

			if actual != tc.expected {
				t.Errorf("got %d; want %d", actual, tc.expected)
			}
		})
	}
}
