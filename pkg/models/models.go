package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicateEmail = errors.New("models: user's email already exists")

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
