package requests

import (
	"encoding/json"
	"fmt"
	"github.com/aliffatulmf/medotcom/parser"
	"io"
	"net/http"
	"strings"
	"time"
)

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

func RequestGET(cookie []parser.Cookie) (*ChatResponse, error) {
	body := bodyBuilder(`{"count": 0}`)

	req, err := NewRequest(http.MethodPost, "https://you.com/api/chat/getUserChats", body, cookie)
	if err != nil {
		return nil, fmt.Errorf("unable to create new request: %w", err)
	}
	defer req.Body.Close()

	res, err := responseClient(req)
	if err != nil {
		return nil, err
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

	if response.Empty() {
		return nil, ErrNoChatFound
	}

	return &response, nil
}

func RequestDELETE(cookie []parser.Cookie, chat *Chat) error {
	body := bodyBuilder(`{"chatId": "%s"}`, chat.ChatID)

	req, err := NewRequest(http.MethodDelete, "https://you.com/api/chat/deleteChat", body, cookie)
	if err != nil {
		return fmt.Errorf("unable to create new request: %w", err)
	}
	defer req.Body.Close()

	res, err := responseClient(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("unexpected status code: %d\n", res.StatusCode)
	}

	return nil
}

func bodyBuilder(format string, args ...interface{}) io.Reader {
	return strings.NewReader(fmt.Sprintf(format, args...))
}

func responseClient(r *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Client: %w", err)
	}

	return res, nil
}
