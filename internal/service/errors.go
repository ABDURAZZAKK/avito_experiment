package service

import "fmt"

var (
	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrNotFound      = fmt.Errorf("not found")
	ErrCannotCreate  = fmt.Errorf("cannot create")
	ErrCannotDelete  = fmt.Errorf("cannot delete")
	ErrCannotGet     = fmt.Errorf("cannot get")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")
	ErrCannotDeleteUser  = fmt.Errorf("cannot delete user")
)
