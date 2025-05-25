package settings

import "mime/multipart"

type SetProfile struct {
	MultipartFile multipart.File
	FileHeader    *multipart.FileHeader
}

type SetTelephone struct {
	Telephone string `json:"telephone"`
}

type UpdateUserInfo struct {
	Username  string     `json:"username"`
	Profile   SetProfile `json:"profile"`
	Email     string     `json:"email"`
	Telephone string     `json:"telephone"`
}
