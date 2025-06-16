package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-coursework/internal/dto/auth"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

type TokenUser struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

var tokenUsers = make(map[string]TokenUser)

func TestSignUp(t *testing.T) {
	var users []auth.SignUpRequest

	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 1; i < 10; i++ {
			if i > 5 {
				var u = auth.SignUpRequest{
					Username:       fmt.Sprintf("User Testing %d", i),
					Email:          fmt.Sprintf("testing%d@gmail.com", i),
					StudyProgramID: 1,
					Password:       fmt.Sprintf("testing%d", i),
					RoleID:         3,
					Batch:          2025,
				}

				mu.Lock()
				users = append(users, u)
				mu.Unlock()
				continue
			}

			var u = auth.SignUpRequest{
				Username:       fmt.Sprintf("User Testing %d", i),
				Email:          fmt.Sprintf("testing%d@gmail.com", i),
				StudyProgramID: 1,
				Password:       fmt.Sprintf("testing%d", i),
				RoleID:         2,
				Batch:          2025,
			}

			mu.Lock()
			users = append(users, u)
			mu.Unlock()
		}
	}()

	wg.Wait()

	for _, user := range users {
		userJson, err := json.Marshal(user)
		if err != nil {
			t.Error(err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/auth/signup", bytes.NewBuffer(userJson))
		if err != nil {
			t.Error(err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}

		var response any
		err = json.Unmarshal(responseBody, &response)
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println("Response:", string(responseBody))
	}
}

var signInResponse = struct {
	Message string `json:"message"`
	Data    struct {
		Email             string  `json:"email"`
		EmailVerified     bool    `json:"email_verified"`
		Telephone         *string `json:"telephone"`
		TelephoneVerified bool    `json:"telephone_verified"`
		Role              string  `json:"role"`
		Semester          string  `json:"semester"`
		Token             string  `json:"token"`
	} `json:"data"`
}{}

func TestSignIn(t *testing.T) {
	var users []auth.SignInRequest

	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 1; i < 10; i++ {
			var u = auth.SignInRequest{
				Email:    fmt.Sprintf("testing%d@gmail.com", i),
				Password: fmt.Sprintf("testing%d", i),
			}

			mu.Lock()
			users = append(users, u)
			mu.Unlock()
		}
	}()

	wg.Wait()

	for i, user := range users {
		userJson, err := json.Marshal(user)
		if err != nil {
			t.Error(err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/auth/signin", bytes.NewBuffer(userJson))
		if err != nil {
			t.Error(err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}

		err = json.Unmarshal(responseBody, &signInResponse)
		if err != nil {
			t.Error(err)
			return
		}

		if signInResponse.Data.Token != "" && signInResponse.Data.Email != "" {
			id := fmt.Sprintf("%d", i+1)
			tokenUsers[id] = TokenUser{
				Email: signInResponse.Data.Email,
				Token: signInResponse.Data.Token,
			}

			fmt.Printf("Successfully saved Email: %s and Token: %s as ID: %s\n", signInResponse.Data.Email, signInResponse.Data.Token, id)

			err := SaveTokensToFile("tokens.json")
			if err != nil {
				t.Errorf("Gagal menyimpan token ke file: %v", err)
			} else {
				fmt.Println("Token berhasil disimpan ke 'tokens.json'")
			}
		} else {
			t.Errorf("No token received for user %s. Response: %s\n", user.Email, string(responseBody))
		}
	}
}

func SaveTokensToFile(filename string) error {
	data, err := json.MarshalIndent(tokenUsers, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func LoadTokenFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &tokenUsers)
}

func TestShowDataTokenUsers(t *testing.T) {
	err := LoadTokenFromFile("tokens.json")
	if err != nil {
		t.Error(err)
		return
	}

	if len(tokenUsers) == 0 {
		fmt.Println("Token users is empty")
	}

	for userID, token := range tokenUsers {
		fmt.Printf("Email: %s | Token: %s\n", userID, token)
	}
}
