package ddbcalc

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
)

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

			gotJSON := string(b)

			f, err := os.ReadFile(path.Join("testdata", fmt.Sprintf("test_%s.json", tt.typ)))
			if err != nil {
				t.Fatal(err)
			}

			wantJSON := string(f)
			if gotJSON != wantJSON {
				t.Errorf("got\n%s\n; want\n%s\n", gotJSON, wantJSON)
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
