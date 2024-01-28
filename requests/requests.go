package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aliffatulmf/medotcom/parser"
)

type RequestOptions struct {
	Payload io.Reader
	Cookies []parser.Cookie
}

func RequestGET(opt *RequestOptions) (*ChatResponse, error) {
	req, err := NewRequest(http.MethodPost, "https://you.com/api/chat/getUserChats", opt.Payload, opt.Cookies)
	if err != nil {
		return nil, fmt.Errorf("unable to create new request: %w", err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	var chatResponse ChatResponse
	if err = json.Unmarshal(data, &chatResponse); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return &chatResponse, nil
}

func RequestDELETE(opt *RequestOptions) error {
	req, err := NewRequest(http.MethodDelete, "https://you.com/api/chat/deleteChat", opt.Payload, opt.Cookies)
	if err != nil {
		return fmt.Errorf("unable to create new request: %w", err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
