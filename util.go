package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
)

type User struct {
	Nome          string `json:"nome"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Administrador string `json:"administrador"`
}

func (u *User) CreateUser(adm bool) {
	u.Nome = gofakeit.Name()
	u.Email = gofakeit.Email()
	u.Password = gofakeit.Password(true, true, true, true, false, 8)
	if adm {
		u.Administrador = "true"
	} else {
		u.Administrador = "false"
	}
}

func SaveData(resp *http.Response, data *map[string]interface{}) {
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Errorf("Aconteceu o erro %e", err)
	}
}

func GenerateID(length int) string {
	var possibleCharacters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, length)
	for i := 0; i < length; i++ {
		b[i] = possibleCharacters[rand.Intn(len(possibleCharacters))]
	}
	return string(b)
}
