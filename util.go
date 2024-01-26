package main

import (
	"github.com/brianvoe/gofakeit/v6"
)

type User struct {
	Nome string `json:"nome"`
	Email string `json:"email"`
	Password string `json:"password"`
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