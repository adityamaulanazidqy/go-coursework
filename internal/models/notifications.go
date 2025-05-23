package models

type Notification struct {
	ID      int    `json:"id" gorm:"primary_key;"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}

func (Notification) TableName() string {
	return "notifications"
}
