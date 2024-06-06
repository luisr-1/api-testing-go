package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	if idValue, ok := result["_id"].(string); ok {
		id = idValue
	} else {
		id = ""
	}

	return req.StatusCode, msg, id, err
}

func putUser(url string, body []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func deleteUser(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
