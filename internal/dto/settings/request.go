package settings

import "mime/multipart"

type SetProfile struct {
	MultipartFile multipart.File
	FileHeader    *multipart.FileHeader
}

type SetTelephone struct {
	Telephone string `json:"telephone"`
}
