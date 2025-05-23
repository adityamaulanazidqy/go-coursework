package constants

type StatusSubmissions struct {
	ID   int    `json:"id" gorm:"primary_key;"`
	Name string `json:"name"`
}

func (StatusSubmissions) TableName() string {
	return "status_submissions"
}
