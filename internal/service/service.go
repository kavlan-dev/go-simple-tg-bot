package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) DogImage(ctx context.Context) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://random.dog/woof.json", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type randomDogResponse struct {
		URL string `json:"url"`
	}

	var dogData randomDogResponse
	err = json.Unmarshal(body, &dogData)
	if err != nil {
		return "", err
	}

	if dogData.URL == "" {
		return "", fmt.Errorf("")
	}

	return dogData.URL, nil
}
