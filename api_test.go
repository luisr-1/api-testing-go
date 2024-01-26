package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func getUsers(url string) (int, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return resp.StatusCode, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, 0, err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return resp.StatusCode, 0, err
	}

	quantity, ok := result["quantidade"].(float64)
	if !ok || quantity < 1 {
		return resp.StatusCode, 0, nil
	}

	users, ok := result["usuarios"].([]interface{})
	if !ok || len(users) != int(quantity) {
		return resp.StatusCode, 0, nil
	}

	return resp.StatusCode, int(quantity), nil
}

func TestGetUsers(t *testing.T) {
	statusCode, quantity, err := getUsers("https://serverest.dev/usuarios")
	if err != nil {
		t.Fatal(err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Era esperado status code 200, recebido %d", statusCode)
	}

	if quantity < 1 {
		t.Errorf("Esperado quantidade de usuarios maior ou igual 1, recebido %d", quantity)
	}
}

func postUser(url string, body []byte) (int, string, string, error) {
	req, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return req.StatusCode, "", "", err
	}
	defer req.Body.Close()
	
	var result map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&result); err != nil {
		return req.StatusCode, fmt.Sprintf("Ocorreu um problema ao decodificar %s", err), "", err
	}

	msg := result["message"].(string)

	var id string
	if idValue, ok := result["_id"].(string); ok  {
		id = idValue
	} else {
		id = ""
	}

	return req.StatusCode, msg, id, err
}
func TestPostUser(t *testing.T) {
	user := &User{}
	user.CreateUser(true)
	userPayload, _ := json.Marshal(user)

	statusCode, msg, id, err := postUser("https://serverest.dev/usuarios", userPayload)
	
	if err != nil {
		t.Fatal(err)
	}

	if statusCode != 201 {
		t.Errorf("O status code esperado era 201, o que foi recebido foi: %d", statusCode)
	}

	if msg != "Cadastro realizado com sucesso" {
		t.Errorf("A mensagem esperada era 'Cadastro realizado com sucesso', a mensagem recebida foi: %s", msg)
	}

	if len(id) != 16 {
		t.Errorf("Era esperado que contivesse um ID com 16 caracteres, o que foi gerado foi com: %d", len(id))
	}
}

func TestPostDuplicatedUser(t *testing.T) { 
	user := &User{}
	user.CreateUser(false)
	userPayload, _ := json.Marshal(user)
	var statusCode int
	var msg string
	var id string
	var err error

	for i := 0; i <= 1; i++ {
		statusCode, msg, id, err = postUser("https://serverest.dev/usuarios", userPayload)
		if err != nil {
			t.Fatal(err)
		}
	}

	if statusCode != 400 {
		t.Errorf("O status code esperado era 400, o que foi recebido foi: %d", statusCode)
	}

	if msg != "Este email já está sendo usado" {
		t.Errorf("Era esperado que a mensagem fosse 'Este email já está sendo usado', mas recebido: %s", msg)
	}

	if id != "" {
		t.Errorf("O esperado era que não houvesse _id gerado, mas foi gerado: %s", id)
	}
}