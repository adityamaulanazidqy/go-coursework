package models

import "time"

type FCMTokens struct {
	ID        int `gorm:"primary_key,AUTO_INCREMENT"`
	UserID    int
	Token     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (FCMTokens) TableName() string {
	return "fcm_tokens"
}
