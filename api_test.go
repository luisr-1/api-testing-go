package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

const URL string = "https://serverest.dev/usuarios/"

func TestGetUsers(t *testing.T) {
	statusCode, quantity, err := getUsers(URL)
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

func TestPostUser(t *testing.T) {
	user := &User{}
	user.CreateUser(true)
	userPayload, _ := json.Marshal(user)

	statusCode, msg, id, err := postUser(URL, userPayload)

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
		statusCode, msg, id, err = postUser(URL, userPayload)
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

// criar usuario com POST, recuperar o id do usuario
// fazer um GET com id, recuperando as informacoes do usuario e guardando para serem comparadas
// modificar os dados do usuario, passar para a requisicao PUT
// fazer assert do status code da request do PUT
// fazer a chamada com o GET id, pegar as informacoes e comparar com as de antes do PUT
func TestPutUser(t *testing.T) {
	user := &User{}
	user.CreateUser(true)

	userPayload, _ := json.Marshal(user)
	_, _, id, err := postUser(URL, userPayload)
	if err != nil {
		t.Errorf("Aconteceu o erro %e", err)
	}

	resp, e := http.Get(URL + id)
	if e != nil {
		t.Errorf("Aconteceu o erro %e", e)
	}

	var dataBefore map[string]interface{}
	SaveData(resp, &dataBefore)

	resp.Body.Close()

	user.Email = gofakeit.Email()
	user.Nome = gofakeit.Name()
	userPayload, _ = json.Marshal(user)

	resp, err = putUser(URL+id, userPayload)
	if err != nil {
		t.Errorf("Error %e", err)
	}

	var dataAfter map[string]interface{}
	SaveData(resp, &dataAfter)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("O status code esperado para essa requisição era 200, encontrado foi %d", resp.StatusCode)
	}

	if dataBefore["email"] == dataAfter["email"] {
		t.Errorf("Houve um erro, o email do usuário era %s e agora era para ser %s", dataBefore["email"], user.Email)
	}

	if dataBefore["nome"] == dataAfter["nome"] {
		t.Errorf("Houve um erro, o nome do usuário era %s e agora era para ser %s", dataBefore["nome"], user.Nome)
	}
}

// Criar usuario, recuperar ID para usar no delete
func TestDeleteUser(t *testing.T) {
	user := User{}
	user.CreateUser(true)

	userPayload, _ := json.Marshal(user)
	_, _, id, err := postUser(URL, userPayload)
	if err != nil {
		t.Errorf("Aconteceu o erro %e", err)
	}

	resp, err := deleteUser(URL + id)
	if err != nil {
		t.Errorf("ocorreu um erro %e", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("O status esperado para a request era: %d, o que foi encontrado foi: %d", http.StatusOK, resp.StatusCode)
	}

	var dataDeleteResponse map[string]interface{}
	SaveData(resp, &dataDeleteResponse)
	expectedMsg := "Registro excluído com sucesso"

	if dataDeleteResponse["message"] != expectedMsg {
		t.Errorf("A mensagem esperada era")
	}

	resp, err = http.Get(URL + id)
	if err != nil {
		t.Errorf("Aconteceu o erro %e", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("O status code esperado para essa request era: %d, o encontrado foi: %d", http.StatusBadRequest, resp.StatusCode)
	}

	var dataGetID map[string]interface{}
	SaveData(resp, &dataGetID)

	expectedMsg = "Usuário não encontrado"

	if dataGetID["message"] != expectedMsg {
		t.Errorf("A mensagem esperada para essa request era: '%s', a recebida foi: '%s'", expectedMsg, dataGetID["message"])
	}
}

func TestDeleteNonExistentUser(t *testing.T) {
	id := GenerateID(16)

	resp, err := deleteUser(URL + id)
	if err != nil {
		t.Errorf("Ocorreu um erro %e", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("O status code esperado era %d e o que foi retornado foi : %d", http.StatusOK, resp.StatusCode)
	}

	var data map[string]interface{}
	SaveData(resp, &data)

	message := "Nenhum registro excluído"
	if data["message"] != message {
		t.Errorf("A mensagem esperada era: '%s'; a obtida foi: %s", message, data["message"])
	}
}
