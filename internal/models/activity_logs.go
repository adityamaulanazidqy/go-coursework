package models

import "time"

type ActivityLog struct {
	ID        int       `json:"id" gorm:"primary_key;"`
	UserID    int       `json:"user_id"`
	Activity  string    `json:"activity"`
	CreatedAt time.Time `json:"created_at"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
