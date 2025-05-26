package models

type UserContactVerification struct {
	ID                int  `json:"id" gorm:"primary_key"`
	UserID            int  `json:"user_id"`
	EmailVerified     bool `json:"email_verified" gorm:"default:false"`
	TelephoneVerified bool `json:"telephone_verified" gorm:"default:false"`
}

func (UserContactVerification) TableName() string {
	return "user_contact_verification"
}
