package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Article struct {
	ID       int
	Title    string
	Content  string
	Category string
}
type User struct {
	ID             int
	FullName       string
	Email          string
	HashedPassword []byte
	Role           string
	Active         bool
	Created        time.Time
}
