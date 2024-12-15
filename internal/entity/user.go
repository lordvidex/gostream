package entity

import (
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

type User struct {
	*gostreamv1.User
}

func (u User) TableName() string {
	return "stream_users"
}

// Key ...
func (u User) Key() uint64 {
	return u.Id
}
