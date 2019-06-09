// +build !integration

package models

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cedric-parisi/fizzbuzz-api/internal/errors"
)

func TestFizzbuzz_Validate(t *testing.T) {
	type fields struct {
		Int1  int
		Int2  int
		Limit int
		Str1  string
		Str2  string
	}
	tests := []struct {
		name   string
		fields fields
		want   []errors.FieldError
	}{
		{
			name: "OK",
			fields: fields{
				3,
				5,
				15,
				"fizz",
				"buzz",
			},
		},
		{
			name: "KO_InvalidArguments",
			fields: fields{
				-1,
				-1,
				-1,
				"",
				"",
			},
			want: []errors.FieldError{
				{
					Name:    "int1",
					Message: "int1 must be greater than 0",
				},
				{
					Name:    "int2",
					Message: "int2 must be greater than 0",
				},
				{
					Name:    "int2",
					Message: "int2 must be greater than int1",
				},
				{
					Name:    "limit",
					Message: "limit must be greater than 0",
				},
				{
					Name:    "str1",
					Message: "str1 must not be empty",
				},
				{
					Name:    "str2",
					Message: "str2 must not be empty",
				},
			},
		},
		{
			name: "KO_LimitTooBig",
			fields: fields{
				3,
				5,
				maxLimit + 1,
				"fizz",
				"buzz",
			},
			want: []errors.FieldError{
				{
					Name:    "limit",
					Message: fmt.Sprintf("limit must be smaller than %d", maxLimit),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fizzbuzz{
				Int1:  tt.fields.Int1,
				Int2:  tt.fields.Int2,
				Limit: tt.fields.Limit,
				Str1:  tt.fields.Str1,
				Str2:  tt.fields.Str2,
			}
			if got := f.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fizzbuzz.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
