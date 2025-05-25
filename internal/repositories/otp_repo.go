package repositories

import (
	"errors"
	"go-coursework/internal/models"
	"gorm.io/gorm"
)

type OtpRepo struct {
	db *gorm.DB
}

func NewOtpRepo(db *gorm.DB) *OtpRepo {
	return &OtpRepo{db: db}
}

func (r *OtpRepo) CheckEmail(email string) (string, error) {
	var user models.Users
	if err := r.db.Where("email = ?", email).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return "Email not found!", err
	}

	return "Email is " + user.Email, nil
}
