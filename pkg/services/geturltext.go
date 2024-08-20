package services

import (
	"errors"
	"io"
	"net/http"
)

type IGetURLText interface {
	GetURLText(url string, client *http.Client) (string, error)
}

type GetURLTextService struct{}

func (s GetURLTextService) GetURLText(url string, client *http.Client) (string, error) {
	if client == nil {
		// Raise error
		return "", errors.New("client is nil")
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/octet-stream")

	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
