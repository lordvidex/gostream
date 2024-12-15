package md5hash

import (
	"testing"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"github.com/stretchr/testify/assert"
)

type sample struct {
	FirstName string `hash:"1,first_name"`
	LastName  string `hash:"0,last_name"`
	embedTop
	Age int // 3 by default
	embedBottom
}

func (s sample) TableName() string {
	return "sample"
}

type embedBottom struct {
	MiddleName string `hash:"1,middle_name"` // After first name
}

type embedTop struct {
	Properties []byte `hash:"2,properties"`
}

func TestValue(t *testing.T) {
	tests := []struct {
		name    string
		it      any
		want    string
		wantErr bool
	}{
		{
			name: "normal struct",
			it:   entity.User{User: &gostreamv1.User{Id: 1, Name: "Test Tester", Age: 10, Nationality: "RU"}},
			want: "1Test Tester10RU",
		},
		{
			name: "normal struct pointer",
			it:   &entity.User{User: &gostreamv1.User{Id: 1, Name: "Test Tester", Age: 10, Nationality: "RU"}},
			want: "1Test Tester10RU",
		},
		{
			name: "raw string value",
			it:   "hello",
			want: "hello",
		},
		{
			name: "struct with tags",
			it: struct {
				Name      string `hash:"2,name"`
				Attention string `hash:"1,attention"`
				Id        int    `hash:"0,id"`
			}{
				Name:      "SomeName",
				Attention: "lowspan",
				Id:        432113213,
			},
			want: "432113213lowspanSomeName",
		},
		{
			name: "embedded struct with clashing indices",
			it: sample{
				FirstName: "First",
				LastName:  "Last",
				Age:       10,
				embedTop: embedTop{
					Properties: []byte(`{"houses":2,"assets":"10"}`),
				},
				embedBottom: embedBottom{MiddleName: "Middle"},
			},
			want: `LastFirstMiddle{"houses":2,"assets":"10"}10`,
		},
		{
			name: "nil pointer",
			it:   (*entity.User)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Value(tt.it)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQuery(t *testing.T) {
	tests := []struct {
		name    string
		it      table
		want    string
		wantErr bool
	}{
		{
			name: "user hash query",
			it:   &entity.User{User: &gostreamv1.User{Id: 1, Name: "Test Tester", Age: 10, Nationality: "RU"}},
			want: `SELECT md5_chain("id"||"name"||"age"||"nationality") FROM stream_users`,
		},
		{
			name: "empty struct",
			it:   entity.User{},
			want: `SELECT md5_chain("id"||"name"||"age"||"nationality") FROM stream_users`,
		},
		{
			name: "sample table query",
			it: sample{
				FirstName: "First",
				LastName:  "Last",
				Age:       10,
				embedTop: embedTop{
					Properties: []byte(`{"houses":2,"assets":"10"}`),
				},
				embedBottom: embedBottom{MiddleName: "Middle"},
			},
			want: `SELECT md5_chain("last_name"||"first_name"||"middle_name"||"properties"||"age") FROM sample`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.it)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChainHashValue(t *testing.T) {
	tests := []struct {
		name    string
		values  []any
		want    string
		wantErr bool
	}{
		{
			name: "empty array",
			want: "",
		},
		{
			name: "single value",
			want: func() string {
				s, _ := HashValue(entity.User{User: &gostreamv1.User{Id: 1, Name: "Tester", Age: 5}})
				return s
			}(),
			values: []any{
				&entity.User{User: &gostreamv1.User{Id: 1, Name: "Tester", Age: 5}},
			},
		},
		{
			name: "multiple values",
			want: "52d67204e243ce24ce2825cb90290467",
			values: []any{
				&entity.User{User: &gostreamv1.User{Id: 1, Name: "Tester", Age: 5}},
				"normal string",
				entity.Pet{Pet: &gostreamv1.Pet{Id: 2, Name: "Bunny", Kind: "pretty"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ChainHashValue(tt.values)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, v)
		})
	}
}
