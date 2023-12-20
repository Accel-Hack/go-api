package repository

import (
	"fmt"
	"time"

	"github.com/Accel-Hack/go-api/internal/domain/sample/model"
	"github.com/google/uuid"
)

type UpdateSampleRow struct {
	ID         string     `xorm:"pk notnull 'ID'"`
	Name       *string    `xorm:"notnull 'NAME'"`
	Birthday   *time.Time `xorm:"notnull 'BIRTHDAY'"`
	IsJapanese *bool      `xorm:"notnull 'IS_JAPANESE'"`
	CreatedAt  *time.Time `xorm:"notnull 'CREATED_AT' created"`
	UpdatedAt  *time.Time `xorm:"notnull 'UPDATED_AT' updated"`
	IsDeleted  *bool      `xorm:"notnull 'IS_DELETED' default 'false'"`
	DeletedAt  *time.Time `xorm:"null 'DELETED_AT'"`
}

type SampleRow struct {
	ID         string    `xorm:"pk notnull 'ID'"`
	Name       string    `xorm:"notnull 'NAME'"`
	Birthday   time.Time `xorm:"notnull 'BIRTHDAY'"`
	IsJapanese bool      `xorm:"notnull 'IS_JAPANESE'"`
	CreatedAt  time.Time `xorm:"notnull 'CREATED_AT' created"`
	UpdatedAt  time.Time `xorm:"notnull 'UPDATED_AT' updated"`
	IsDeleted  bool      `xorm:"notnull 'IS_DELETED' default 'false'"`
	DeletedAt  time.Time `xorm:"null 'DELETED_AT'"`
}

func (r SampleRow) toSample() (*model.Sample, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	return &model.Sample{
		ID:         id,
		Name:       r.Name,
		Birthday:   r.Birthday,
		IsJapanese: r.IsJapanese,
	}, err
}
