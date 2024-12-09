package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/antoniofmoliveira/courses/dto"
)

func main() {
	user := dto.GetJWTInput{
		Email:    "user@gmail.com",
		Password: "123456",
	}

	jsonbytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:8081/users/generate_token", bytes.NewBuffer(jsonbytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
