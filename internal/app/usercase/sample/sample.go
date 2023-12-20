package sample

import (
	"context"
	"time"

	"github.com/Accel-Hack/go-api/internal/domain/sample/model"
	"github.com/google/uuid"
)

const (
	DefaultLimit  = 20
	DefaultOffset = 0
)

type UpdateQuery struct {
	ID         uuid.UUID
	Name       *string
	Birthday   *time.Time
	IsJapanese *bool
}

type SampleRepository interface {
	// SELECT `ID`, `NAME` , `BIRTHDAY`, `IS_JAPANESE` FROM @@tablew WHERE `IS_DELETED` IS FALSE AND `ID` = @id
	FindByID(ctx context.Context, id uuid.UUID) (*model.Sample, error)
	// SELECT `ID`, `NAME` , `BIRTHDAY`, `IS_JAPANESE`, COUNT(*) OVER () AS TOTAL FROM @@table WHERE `IS_DELETED` IS FALSE {{if name != ""}} AND `NAME` LIKE concat("%",@name,"%") {{end}} LIMIT @limit OFFSET @offset
	FindByNameLike(ctx context.Context, name string, offset, limit int) (*model.PagedSamples, error)
	// INSERT INTO @@table (`ID`, `NAME`, `BIRTHDAY`, `IS_JAPANESE`) VALUES (@sample.id, @sample.name, @sample.birthday, @sample.isJapanese) ON DUPLICATE KEY UPDATE `NAME` = @sample.name, `BIRTHDAY` = @sample.birthday, `IS_JAPANESE` = @sample.isJapanese
	Insert(ctx context.Context, sample *model.Sample) error
	// UPDATE @@table SET `NAME` = @sample.Name, `BIRTHDAY` = @sample.Birthday, `IS_JAPANESE` = @sample.isJapanese WHERE `ID` = @sample.ID
	Update(ctx context.Context, sample UpdateQuery) error
	// UPDATE @@table SET `IS_DELETED` = true, `DELETED_AT` = CURRENT_TIMESTAMP WHERE `ID` = @id
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	Repository SampleRepository
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*model.Sample, error) {
	return u.Repository.FindByID(ctx, id)
}

func (u *Usecase) Search(ctx context.Context, name string, limit, offset *int) (*model.PagedSamples, error) {
	l := DefaultLimit
	if limit != nil {
		l = *limit
	}
	o := DefaultOffset
	if offset != nil {
		o = *offset
	}
	return u.Repository.FindByNameLike(ctx, name, o, l)
}

type AddQuery struct {
	Name       string
	Birthday   time.Time
	IsJapanese bool
}

func (u *Usecase) Add(ctx context.Context, q AddQuery) (uuid.UUID, error) {
	sample := model.NewSample(q.Name, q.Birthday, q.IsJapanese)
	err := u.Repository.Insert(ctx, sample)
	if err != nil {
		return uuid.UUID{}, err
	}
	return sample.ID, nil
}

func (u *Usecase) Edit(ctx context.Context, q UpdateQuery) error {
	return u.Repository.Update(ctx, UpdateQuery{
		ID:         q.ID,
		Name:       q.Name,
		Birthday:   q.Birthday,
		IsJapanese: q.IsJapanese,
	})
}

func (u *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.Repository.DeleteByID(ctx, id)
}
