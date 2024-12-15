package entity

import (
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

type Pet struct {
	*gostreamv1.Pet
}

// TableName implements md5hash.table
func (u Pet) TableName() string {
	return "pets"
}

// Key ...
func (u Pet) Key() uint64 {
	return u.Id
}
