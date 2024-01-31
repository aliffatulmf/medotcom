package requests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aliffatulmf/medotcom/parser"
)

type Chat struct {
	ID         string `json:"_id"`
	ChatID     string `json:"chat_id"`
	Title      string `json:"title"`
	DateUpdate string `json:"date_updated"`
}

type ChatResponse struct {
	Chats []Chat `json:"chats"`
}

type RequestOptions struct {
	Payload io.Reader
	Cookies []parser.Cookie
}

func NewRequest(method, url string, body io.Reader, cookies []parser.Cookie) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %s", err)
	}

	req.Header.Set("Host", "you.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.0.0")

	for _, c := range cookies {
		if err = c.HttpCookie().Valid(); err != nil {
			return nil, fmt.Errorf("cookie not valid: %s", err)
		}
		req.Header.Add("Cookie", c.HttpCookie().String())
	}

	return req, err
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

	var response ChatResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response: %w", err)
	}

	return &response, nil
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
