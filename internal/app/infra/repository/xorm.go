package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Accel-Hack/go-api/internal/app/usercase/sample"
	"github.com/Accel-Hack/go-api/internal/domain/sample/model"
	"github.com/google/uuid"
	"xorm.io/xorm"
)

type SampleXorm struct {
	e     *xorm.Engine
	table string
}

func NewSampleXorm(e *xorm.Engine, table string) *SampleXorm {
	return &SampleXorm{
		e:     e,
		table: table,
	}
}

var ErrNotFound = errors.New("not found")

// FindByID implements sample.SampleRepository.
func (r *SampleXorm) FindByID(ctx context.Context, id uuid.UUID) (*model.Sample, error) {
	sampleRow := SampleRow{}
	ok, err := r.e.Context(ctx).Table(r.table).ID(id.String()).Where("IS_DELETED = ?", false).Get(&sampleRow)
	if !ok {
		log.Printf("ID(%q) not found", id.String())
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return sampleRow.toSample()
}

// FindByNameLike implements sample.SampleRepository.
func (r *SampleXorm) FindByNameLike(ctx context.Context, name string, offset int, limit int) (*model.PagedSamples, error) {
	sampleRows := []SampleRow{}
	count, err := r.e.Context(ctx).Table(r.table).
		Where("IS_DELETED = ?", false).
		Where("name LIKE ?", contains(name)).
		Limit(limit, offset).
		FindAndCount(&sampleRows)
	if err != nil {
		return nil, err
	}
	samples := make([]model.Sample, len(sampleRows))
	for i, s := range sampleRows {
		sample, err := s.toSample()
		if err != nil {
			return nil, err
		}
		samples[i] = *sample
	}
	return model.NewPagedSamples(int(count), samples)
}

// Insert implements sample.SampleRepository.
func (r *SampleXorm) Insert(ctx context.Context, s *model.Sample) error {
	newRow := SampleRow{
		ID:         s.ID.String(),
		Name:       s.Name,
		Birthday:   s.Birthday,
		IsJapanese: s.IsJapanese,
	}
	_, err := r.e.Context(ctx).Table(r.table).Insert(&newRow)
	if err != nil {
		return fmt.Errorf("insert %#v: %w", newRow, err)
	}
	return nil
}

// Update implements sample.SampleRepository.
func (r *SampleXorm) Update(ctx context.Context, query sample.UpdateQuery) error {
	updateRow := UpdateSampleRow{
		ID:         query.ID.String(),
		Name:       query.Name,
		Birthday:   query.Birthday,
		IsJapanese: query.IsJapanese,
	}
	log.Printf("[DEBUG] update sample with %#v", updateRow)
	if _, err := r.e.Context(ctx).Table(r.table).ID(query.ID.String()).Update(&updateRow); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

// DeleteByID implements sample.SampleRepository.
func (r *SampleXorm) DeleteByID(ctx context.Context, id uuid.UUID) error {
	var (
		isDeleted = true
		now       = time.Now()
	)
	_, err := r.e.Context(ctx).Table(r.table).ID(id.String()).Update(&UpdateSampleRow{
		ID:        id.String(),
		IsDeleted: &isDeleted,
		DeletedAt: &now,
	})
	if err != nil {
		return fmt.Errorf("delete %s: %w", id, err)
	}
	return nil
}

var _ (sample.SampleRepository) = (*SampleXorm)(nil)
