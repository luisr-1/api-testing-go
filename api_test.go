package main

import (
	"encoding/json"
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

func TestUserEndpoint(t *testing.T) {
	statusCode, quantity, err := getUsers("https://serverest.dev/usuarios/")
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
