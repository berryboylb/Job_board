package files

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

const baseUrl = "https://api.uploadfly.cloud"

type Request struct {
	APIKey string
}

func NewAPIClient(apiKey string) *Request {
	return &Request{
		APIKey: apiKey,
	}
}

type UploadRequest struct {
	File             multipart.FileHeader `json:"file"`
	Filename         string               `json:"filename"`
	MaxFileSize      string               `json:"maxFileSize"`
	AllowedFileTypes string               `json:"allowedFileTypes"`
}

// Post sends a POST request to the specified endpoint with the given payload.
func (c *Request) Post(payload UploadRequest) (*http.Response, error) {
	url := baseUrl + "/upload"

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := http.DefaultClient
	return client.Do(req)
}

type DeleteRequest struct {
	FileUrl string `json:"file_url"`
}

// Post sends a POST request to the specified endpoint with the given payload.
func (c *Request) Delete(payload DeleteRequest) (*http.Response, error) {
	url := baseUrl + "/delete"

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := http.DefaultClient
	return client.Do(req)
}
