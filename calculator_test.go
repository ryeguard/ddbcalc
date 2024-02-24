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
		typ  string
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
		{
			name: "string",
			typ:  "string",
			item: struct{ StringField string }{StringField: "101"},
			want: 11 + 3,
		},
		{
			name: "int",
			typ:  "int",
			item: struct{ IntField int }{IntField: 101},
			want: 8 + 3,
		},
		{
			name: "float64",
			typ:  "float64",
			item: struct{ Float64Field float64 }{Float64Field: 2.1},
			want: 12 + 3,
		},
		{
			name: "strings",
			typ:  "[]string",
			item: struct{ StringsField []string }{StringsField: []string{"a", "b", "c"}},
			want: 3 + 12 + 3*(1+1),
		},
		{
			name: "bool",
			typ:  "bool",
			item: struct{ BoolField bool }{BoolField: true},
			want: 9 + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// skip json verification if typ is not provided
			if tt.typ != "" {
				b, err := json.MarshalIndent(tt.item, "", "\t")
				if err != nil {
					t.Fatal(err)
				}

				gotJson := string(b)

				f, err := os.ReadFile(path.Join("testdata", fmt.Sprintf("test_%s.json", tt.name)))
				if err != nil {
					t.Fatal(err)
				}

				wantJson := string(f)
				if gotJson != wantJson {
					t.Errorf("got\n%s\n; want\n%s\n", gotJson, wantJson)
				}
			}

			got := StructSizeInBytes(tt.item)
			if got != tt.want {
				t.Errorf("got %d; want %d", got, tt.want)
			}
		})
	}
}
