package inmemory

import (
	"context"
	"testing"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	arr := NewArray[uint64, *entity.User]()
	ch, err := NewCache(context.Background(), arr, nil)
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

func TestArray_Get(t *testing.T) {
	arr := NewArray[uint64, *entity.User]()
	_, ok := arr.Get(0)
	assert.Equal(t, false, ok)
}

func TestArray_Add(t *testing.T){
	var arr Array[uint64, *entity.User]
	arr.add(1, &entity.User{})

}