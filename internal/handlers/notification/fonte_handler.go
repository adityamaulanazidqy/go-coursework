package notification

import (
	"bytes"
	"errors"
	"github.com/go-jose/go-jose/v4/json"
	"go-coursework/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type FonteHandler struct {
	db *gorm.DB
}

func NewFonteHandler(db *gorm.DB) *FonteHandler {
	return &FonteHandler{
		db: db,
	}
}

func (h *FonteHandler) SendMessage(n, t int, m string) error {
	var (
		url   = "https://api.fonnte.com/send"
		token = t

		payload = map[string]interface{}{
			"target":  n,
			"message": m,
		}
	)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", strconv.Itoa(token))
	req.Header.Set("Content-Type", "application/json")

	return nil
}

func (h *FonteHandler) GetNumber(i int) (
	n int,
	err error,
) {
	var u models.Users
	if err := h.db.Where("id = ?", i).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("user not found")
			return 0, err
		}

		return 0, err
	}

	if *u.Telephone == "" {
		err := errors.New("telephone number not found")
		return 0, err
	}

	return n, nil
}
