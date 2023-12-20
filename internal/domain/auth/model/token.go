package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	id                  uuid.UUID `json:"-"`
	accessToken         string
	encryptRefreshToken string
	expiresAt           time.Time
}
