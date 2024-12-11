package entity

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

type User struct {
	*gostreamv1.User
}

// Key ...
func (u User) Key() uint64 {
	return u.Id
}

// Hash ...
func (u User) Hash() string {
	str := u.UniqueString()
	h := md5.New()
	io.WriteString(h, str)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (u User) UniqueString() string {
	var buf strings.Builder

	fmt.Fprintf(&buf, "%d", u.Id)
	fmt.Fprintf(&buf, "%s", u.Name)
	fmt.Fprintf(&buf, "%d", u.Age)
	fmt.Fprintf(&buf, "%s", u.Nationality)
	return buf.String()
}
