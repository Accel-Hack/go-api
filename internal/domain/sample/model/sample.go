package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Sample struct {
	ID         uuid.UUID
	Name       string
	Birthday   time.Time
	IsJapanese bool
}

func NewSample(name string, birthday time.Time, isJapanese bool) *Sample {
	return &Sample{
		ID:         uuid.New(),
		Name:       name,
		Birthday:   birthday,
		IsJapanese: isJapanese,
	}
}

type PagedSamples struct {
	Total   int
	Samples []Sample
}

var ErrNegativeTotal = errors.New("total count is negative")

func NewPagedSamples(total int, samples []Sample) (*PagedSamples, error) {
	if total < 0 {
		return nil, ErrNegativeTotal
	}
	return &PagedSamples{
		Total:   total,
		Samples: samples,
	}, nil
}
