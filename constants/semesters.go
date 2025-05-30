package constants

type Semesters struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

func (Semesters) TableName() string {
	return "semesters"
}
