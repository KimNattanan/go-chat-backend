package persistent

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	db := r.db.WithContext(ctx)
	return db.Create(user).Error
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	db := r.db.WithContext(ctx)
	var user entity.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := r.db.WithContext(ctx)
	var user entity.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	db := r.db.WithContext(ctx)
	result := db.Delete(&entity.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
