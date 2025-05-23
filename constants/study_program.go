package constants

type StudyPrograms struct {
	ID   int    `json:"id" gorm:"primary_key;"`
	Name string `json:"name"`
}

func (StudyPrograms) TableName() string {
	return "study_programs"
}
