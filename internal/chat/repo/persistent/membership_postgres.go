package persistent

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"gorm.io/gorm"
)

type MembershipRepo struct {
	db *gorm.DB
}

func NewMembershipRepo(db *gorm.DB) *MembershipRepo {
	return &MembershipRepo{
		db: db,
	}
}

func (r *MembershipRepo) Create(ctx context.Context, membership *entity.Membership) error {
	db := r.db.WithContext(ctx)
	return db.Create(membership).Error
}

func (r *MembershipRepo) FindByRoomID(ctx context.Context, roomID string) ([]*entity.Membership, error) {
	db := r.db.WithContext(ctx)
	var membershipValues []entity.Membership
	if err := db.Where("room_id = ?", roomID).Find(&membershipValues).Error; err != nil {
		return nil, err
	}

	memberships := make([]*entity.Membership, len(membershipValues))
	for i := range membershipValues {
		memberships[i] = &membershipValues[i]
	}
	return memberships, nil
}

func (r *MembershipRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.Membership, error) {
	db := r.db.WithContext(ctx)
	var membershipValues []entity.Membership
	if err := db.Where("user_id = ?", userID).Find(&membershipValues).Error; err != nil {
		return nil, err
	}

	memberships := make([]*entity.Membership, len(membershipValues))
	for i := range membershipValues {
		memberships[i] = &membershipValues[i]
	}
	return memberships, nil
}

func (r *MembershipRepo) FindByRoomIDAndUserID(ctx context.Context, roomID, userID string) (*entity.Membership, error) {
	db := r.db.WithContext(ctx)
	var membership entity.Membership
	if err := db.Where("room_id = ?", roomID).Where("user_id = ?", userID).First(&membership).Error; err != nil {
		return &entity.Membership{}, err
	}
	return &membership, nil
}

func (r *MembershipRepo) Delete(ctx context.Context, roomID, userID string) error {
	db := r.db.WithContext(ctx)
	result := db.Where("room_id = ?", roomID).Where("user_id = ?", userID).Delete(&entity.Membership{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MembershipRepo) DeleteByUserID(ctx context.Context, userID string) error {
	db := r.db.WithContext(ctx)
	result := db.Delete(&entity.Membership{}, "user_id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
