package persistent

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"gorm.io/gorm"
)

type ProfileRepo struct {
	db *gorm.DB
}

func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{
		db: db,
	}
}

func (r *ProfileRepo) Create(ctx context.Context, profile *entity.Profile) error {
	db := r.db.WithContext(ctx)
	return db.Create(profile).Error
}

func (r *ProfileRepo) FindByUserID(ctx context.Context, userID string) (*entity.Profile, error) {
	db := r.db.WithContext(ctx)
	var profile entity.Profile
	if err := db.First(&profile, "user_id = ?", userID).Error; err != nil {
		return &entity.Profile{}, err
	}
	return &profile, nil
}

func (r *ProfileRepo) Patch(ctx context.Context, userID string, profile *entity.Profile) error {
	db := r.db.WithContext(ctx)
	result := db.Model(&entity.Profile{}).Where("user_id = ?", userID).Updates(map[string]string{
		"name": profile.Name,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ProfileRepo) Delete(ctx context.Context, userID string) error {
	db := r.db.WithContext(ctx)
	result := db.Delete(&entity.Profile{}, "user_id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
