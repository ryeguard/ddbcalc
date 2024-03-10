package ddbcalc

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestSizeInBytesNil(t *testing.T) {
	actual, err := sizeInBytes(nil)
	if err != nil {
		t.Fatal(err)
	}
	if actual != 0 {
		t.Errorf("got %d; want 0", actual)
	}
}

func TestSizeInBytes(t *testing.T) {
	var tests = []struct {
		name     string
		item     types.AttributeValue
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
			name: "list",
			item: &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "abc"},
				&types.AttributeValueMemberN{Value: "123"},
				&types.AttributeValueMemberB{Value: []byte{1, 2, 3}},
			}},
			expected: 3 + // list overhead
				3*(3+1), // 3 elements 3 bytes each + 1 byte overhead
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			actual, err := sizeInBytes(&tt.item)
			if err != nil {
				t.Fatal(err)
			}
			if actual != tt.expected {
				t.Errorf("got %d; want %d", actual, tt.expected)
			}
		})
	}
}

func TestStructSizeInBytes(t *testing.T) {
	var tests = []struct {
		name string
		item interface{}
		want int
	}{
		{
			name: "unexported fields are ignored",
			item: struct{ unexported string }{unexported: "abc"},
			want: 0,
		},
		{
			name: `fields tagged with "dynamodbav:"-" are ignored`,
			item: struct {
				ExportedTagged string `dynamodbav:"-"`
			}{ExportedTagged: "abc"},
			want: 0,
		},
		{
			name: "json tag is ignored",
			item: struct {
				LongFieldName string `json:"1"`
			}{LongFieldName: "22"},
			want: 13 + 2,
		},
		{
			name: "dynamodbav tag is used",
			item: struct {
				LongFieldName string `dynamodbav:"1"`
			}{LongFieldName: "22"},
			want: 1 + 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StructSizeInBytes(tt.item)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %d; want %d", got, tt.want)
			}
		})
	}
}

func TestStrctSizeInBytesOfTypes(t *testing.T) {
	var tests = []struct {
		typ  string
		item interface{}
		want int
	}{
		{
			typ:  "string",
			item: struct{ StringField string }{StringField: "101"},
			want: 11 + 3,
		},
		{
			typ:  "int",
			item: struct{ IntField int }{IntField: 101},
			want: 8 + 3,
		},
		{
			typ:  "float64",
			item: struct{ Float64Field float64 }{Float64Field: 2.1},
			want: 12 + 3,
		},
		{
			typ:  "[]string",
			item: struct{ StringsField []string }{StringsField: []string{"a", "b", "c"}},
			want: 3 + 12 + 3*(1+1),
		},
		{
			typ:  "bool",
			item: struct{ BoolField bool }{BoolField: true},
			want: 9 + 1,
		},
		{
			typ:  "byte",
			item: struct{ ByteField byte }{ByteField: 123},
			want: 9 + 3,
		},
		{
			// []byte will be marshaled as Binary data (B)
			typ:  "[]byte",
			item: struct{ BytesField []byte }{BytesField: []byte{1, 2, 3}},
			want: 13,
		},
		{
			// [][]byte will be marshaled as Binary Set data (BS)
			typ: "[][]byte",
			item: struct{ BytesField [][]byte }{
				BytesField: [][]byte{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			want: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			b, err := json.MarshalIndent(tt.item, "", "\t")
			if err != nil {
				t.Fatal(err)
			}

			gotJson := string(b)

			f, err := os.ReadFile(path.Join("testdata", fmt.Sprintf("test_%s.json", tt.typ)))
			if err != nil {
				t.Fatal(err)
			}

			wantJson := string(f)
			if gotJson != wantJson {
				t.Errorf("got\n%s\n; want\n%s\n", gotJson, wantJson)
			}

			got, err := StructSizeInBytes(tt.item)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %d; want %d", got, tt.want)
			}
		})
	}
}
