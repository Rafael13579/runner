package invoker

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SignResponse struct {
	Signature string `json:"signature"`
	Valid     bool   `json:"valid"`
	Message   string `json:"message"`
}

func Sign(content string) (string, error) {

	body := map[string]string{
		"content": content,
		"token":   "1234",
	}

	jsonData, _ := json.Marshal(body)

	resp, err := http.Post("http://localhost:8080/api/sign", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result SignResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Signature, nil
}

func Validate(content, signature string) (bool, error) {

	body := map[string]string{
		"content":   content,
		"signature": signature,
	}

	jsonData, _ := json.Marshal(body)

	resp, err := http.Post("http://localhost:8080/api/validate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result SignResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Valid, nil
}
