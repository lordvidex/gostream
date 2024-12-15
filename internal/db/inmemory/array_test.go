package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

func TestArray_computeHash(t *testing.T) {
	arr := NewArray[uint64, *entity.User]()
	ch, err := NewCache(context.Background(), arr)
	arr.AddAll([]*entity.User{
		{User: &gostreamv1.User{Id: 2, Nationality: "RU", Age: 22, Name: "Evans"}},
		{User: &gostreamv1.User{Id: 1, Nationality: "CA", Age: 4, Name: "Joe"}},
		{User: &gostreamv1.User{Id: 3, Nationality: "RU", Age: 14, Name: "Evans Jones"}},
	})
	require.NoError(t, err)
	ch.computeHash()
	fromDB := "37a3282bedee25ce3ad52bc24e08c6de" // select md5_chain(id || name || age || nationality) FROM stream_users
	assert.Equal(t, fromDB, ch.dataHash)
}

func TestArray(t *testing.T) {
	arr := NewArray[uint64, entity.User]()

	// add items
	arr.Add(entity.User{User: &gostreamv1.User{Id: 1}})
	arr.Add(entity.User{User: &gostreamv1.User{Id: 2}})
	arr.Add(entity.User{User: &gostreamv1.User{Id: 3}})
	arr.Add(entity.User{User: &gostreamv1.User{Id: 4}})
	arr.Add(entity.User{User: &gostreamv1.User{Id: 5}})
	arr.Add(entity.User{User: &gostreamv1.User{Id: 6}})

	// get items
	for i := range 5 {
		_, ok := arr.Get(uint64(i + 1))
		assert.Truef(t, ok, "result for key: %d", i)
	}
	_, ok := arr.Get(7)
	assert.False(t, ok)

	// delete items
	arr.Remove(4)

	// get items
	for i := range 5 {
		key := i + 1
		_, ok := arr.Get(uint64(key))
		if key == 4 {
			assert.False(t, ok)
		} else {
			assert.True(t, ok)
		}
	}

	// clear
	arr.Clear()

	// iter should return immediately
	for range arr.Iter() {
		t.Fail()
	}
}

func TestArray_Add(t *testing.T) {
	tests := []struct {
		newItem entity.User
		name    string
		initial []entity.User
	}{
		{
			name: "add already existing item",
			initial: []entity.User{
				{User: &gostreamv1.User{Id: 1}},
				{User: &gostreamv1.User{Id: 2}},
				{User: &gostreamv1.User{Id: 3}},
				{User: &gostreamv1.User{Id: 4}},
				{User: &gostreamv1.User{Id: 5}},
				{User: &gostreamv1.User{Id: 6}},
			},
			newItem: entity.User{User: &gostreamv1.User{Id: 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := NewArray[uint64, entity.User]()
			arr.AddAll(tt.initial)

			arr.Add(tt.newItem)

			v, ok := arr.Get(tt.newItem.Key())
			assert.True(t, ok)
			assert.Equal(t, v, tt.newItem)
		})
	}
}
