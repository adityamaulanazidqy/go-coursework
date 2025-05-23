package constants

type Roles struct {
	ID   int    `json:"id" gorm:"primary_key;"`
	Name string `json:"name"`
}

func (Roles) TableName() string {
	return "roles"
}
