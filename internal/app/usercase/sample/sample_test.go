package sample

import (
	"context"
	"testing"

	"github.com/Accel-Hack/go-api/internal/domain/sample/model"
	"github.com/google/uuid"
)

func TestUsecase_Get(t *testing.T) {
}

type stubRepository struct{}

// DeleteByID implements SampleRepository.
func (*stubRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// FindByID implements SampleRepository.
func (*stubRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Sample, error) {
	panic("unimplemented")
}

// FindByNameLike implements SampleRepository.
func (*stubRepository) FindByNameLike(ctx context.Context, name string, offset int, limit int) (*model.PagedSamples, error) {
	panic("unimplemented")
}

// Insert implements SampleRepository.
func (*stubRepository) Insert(ctx context.Context, sample *model.Sample) error {
	panic("unimplemented")
}

// Update implements SampleRepository.
func (*stubRepository) Update(ctx context.Context, sample UpdateQuery) error {
	panic("unimplemented")
}

var _ SampleRepository = (*stubRepository)(nil)
