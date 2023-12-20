package service

import (
	"time"

	"github.com/Accel-Hack/go-api/internal/domain/sample/model"
)

func UpdateSample(old *model.Sample, name *string, birthday *time.Time, isJapanese *bool) *model.Sample {
	newName := old.Name
	if name != nil {
		newName = *name
	}
	newBirthday := old.Birthday
	if birthday != nil {
		newBirthday = *birthday
	}
	newIsJapanese := old.IsJapanese
	if isJapanese != nil {
		newIsJapanese = *isJapanese
	}
	return &model.Sample{
		ID:         old.ID,
		Name:       newName,
		Birthday:   newBirthday,
		IsJapanese: newIsJapanese,
	}
}
