package calc

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestSizeOfNil(t *testing.T) {
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

func TestSizeOfBasicTypes(t *testing.T) {
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
			name:     "S (string) 1",
			item:     &types.AttributeValueMemberS{Value: "abc"},
			expected: 3,
		},
		{
			name:     "S (string) 2",
			item:     &types.AttributeValueMemberS{Value: "12345"},
			expected: 5,
		},
		{
			name:     "N (number) 1",
			item:     &types.AttributeValueMemberN{Value: "123"},
			expected: 3,
		},
		{
			name:     "N (number) 2",
			item:     &types.AttributeValueMemberN{Value: "123.4"},
			expected: 5,
		},
		{
			name:     "B (binary) 1",
			item:     &types.AttributeValueMemberB{Value: []byte{1, 2, 3}},
			expected: 3,
		},
		{
			name:     "B (binary) 2",
			item:     &types.AttributeValueMemberB{Value: []byte{1, 2, 3, 4, 255}},
			expected: 5,
		},
		{
			name:     "BOOL (bool) 1",
			item:     &types.AttributeValueMemberBOOL{Value: false},
			expected: 1,
		},
		{
			name:     "BOOL (bool) 2",
			item:     &types.AttributeValueMemberBOOL{Value: true},
			expected: 1,
		},
		{
			name:     "NULL (null) 1",
			item:     &types.AttributeValueMemberNULL{Value: false},
			expected: 1,
		},
		{
			name:     "NULL (null) 2",
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

func TestSizeOfSet(t *testing.T) {
	var tests = []struct {
		item     types.AttributeValue
		name     string
		expected int
	}{
		{
			name: "BS (binary set) 1",
			item: &types.AttributeValueMemberBS{Value: [][]byte{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			}},
			expected: 9,
		},
		{
			name: "BS (binary set) 2",
			item: &types.AttributeValueMemberBS{Value: [][]byte{
				{1, 253},
				{2, 254},
				{3, 255},
			}},
			expected: 6,
		},
		{
			name:     "NS (number set) 1",
			item:     &types.AttributeValueMemberNS{Value: []string{"1", "2", "3"}},
			expected: 3,
		},
		{
			name:     "NS (number set) 2",
			item:     &types.AttributeValueMemberNS{Value: []string{"1", "2", "3", "4.5"}},
			expected: 6,
		},
		{
			name:     "SS (string set) 1",
			item:     &types.AttributeValueMemberSS{Value: []string{"a", "b", "c"}},
			expected: 3,
		},
		{
			name:     "SS (string set) 2",
			item:     &types.AttributeValueMemberSS{Value: []string{"a", "b", "c", "d"}},
			expected: 4,
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

func TestSizeOfList(t *testing.T) {
	var tests = []struct {
		item     types.AttributeValue
		name     string
		expected int
	}{
		{
			name: "L (list) - empty",
			item: &types.AttributeValueMemberL{
				Value: []types.AttributeValue{},
			},
			expected: 3,
		},
		{
			name: "L (list) - single BOOL value",
			item: &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberBOOL{Value: false},
				},
			},
			expected: 5,
		},
		{
			name: "L (list) - two BOOL values",
			item: &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberBOOL{Value: false},
					&types.AttributeValueMemberBOOL{Value: true},
				},
			},
			expected: 7,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SizeInBytes(&tc.item)

			if actual != tc.expected {
				t.Errorf("got %d, want %d", actual, tc.expected)
			}
		})
	}
}

func TestSizeOfMap(t *testing.T) {
	var tests = []struct {
		item     types.AttributeValue
		name     string
		expected int
	}{
		{
			name:     "M (map) - empty",
			item:     &types.AttributeValueMemberM{},
			expected: 3,
		},
		{
			name: "M (map) - empty",
			item: &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{},
			},
			expected: 3,
		},
		{
			name: "M (map) - single BOOL value",
			item: &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"key1": &types.AttributeValueMemberBOOL{},
				},
			},
			expected: 9,
		},
		{
			name: "M (map) - single BOOL value",
			item: &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"key1": &types.AttributeValueMemberBOOL{Value: false},
				},
			},
			expected: 9,
		},
		{
			name: "M (map) - two BOOL values",
			item: &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"key1": &types.AttributeValueMemberBOOL{Value: false},
					"key2": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
			expected: 15,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := SizeInBytes(&tc.item)

			if actual != tc.expected {
				t.Errorf("got %d, want %d", actual, tc.expected)
			}
		})
	}

}
