package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate stringer -type=Actor
type Actor int

const (
	ActorUnknown Actor = iota // UNKNOWN
	ActorSystem               // SYSTEM
	ActorManager              // MANAGER
	ActorUser                 // USER
)

type User struct {
	id              uuid.UUID
	userName        string
	encryptPassword string
	actor           Actor
	resetCode       *string
	resetUntil      *time.Time
	tokens          []Token
}

var (
	ErrEmptyUserNameEmpty       = errors.New("user name is empty")
	ErrEmptyEncryptRefreshToken = errors.New("encrypt password is empty")
)

func NewUser(id uuid.UUID, userName, encryptPassword string, actor Actor, resetCode *string, resetUntil *time.Time, tokens []Token) (*User, error) {
	if userName == "" {
		return nil, ErrEmptyUserNameEmpty
	}
	if encryptPassword == "" {
		return nil, ErrEmptyEncryptRefreshToken
	}
	if tokens == nil {
		tokens = []Token{}
	}
	return &User{
		id:              id,
		userName:        userName,
		encryptPassword: encryptPassword,
		actor:           actor,
		resetCode:       resetCode,
		resetUntil:      resetUntil,
		tokens:          tokens,
	}, nil
}
