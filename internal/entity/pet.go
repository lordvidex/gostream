package entity

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

type Pet struct {
	*gostreamv1.Pet
}

// Key ...
func (u Pet) Key() uint64 {
	return u.Id
}

// Hash ...
func (u Pet) Hash() string {
	str := u.UniqueString()
	h := md5.New()
	io.WriteString(h, str)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (u Pet) UniqueString() string {
	var buf strings.Builder

	fmt.Fprintf(&buf, "%d", u.Id)
	fmt.Fprintf(&buf, "%s", u.Name)
	fmt.Fprintf(&buf, "%s", u.Kind)
	fmt.Fprintf(&buf, "%d", u.Age)
	return buf.String()
}
