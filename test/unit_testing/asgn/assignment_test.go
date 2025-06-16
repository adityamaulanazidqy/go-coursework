package asgn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

type TokenUser struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

var tokenUsers = make(map[string]TokenUser)

func TestPostAssignment(t *testing.T) {
	err := LoadTokenFromFile("../users/tokens.json")
	if err != nil {
		t.Fatalf("Gagal membuka file json: %v", err)
		return
	}

	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("%d", i)
		token, ok := tokenUsers[id]
		if !ok {
			t.Logf("Token untuk user %s tidak ditemukan, dilewati", id)
			continue
		}

		file, err := os.Open("../../../assets/images_test/e-commers-go.jpg")
		if err != nil {
			t.Fatalf("Gagal membuka file: %v", err)
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		writer.WriteField("title", fmt.Sprintf("Testing %d", i))
		writer.WriteField("description", "In Go, middleware functions act as interceptors for HTTP requests.")
		writer.WriteField("deadline", "2025-06-12 23:59:59")
		writer.WriteField("semester_id", "1")
		writer.WriteField("study_program_id", "1")

		part, err := writer.CreateFormFile("file", "e-commers-go.jpg")
		if err != nil {
			t.Fatalf("Gagal membuat form file: %v", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("Gagal menyalin file ke form: %v", err)
		}

		writer.Close()

		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/assignments", body)
		if err != nil {
			t.Fatalf("Gagal membuat request: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token.Token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Gagal mengirim request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			t.Errorf("Dosen %s gagal: %s, response: %s", id, resp.Status, string(respBody))
		} else {
			t.Logf("Dosen %s berhasil: %s", id, resp.Status)
		}
	}
}

func LoadTokenFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &tokenUsers)
}
