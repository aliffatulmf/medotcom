package requests

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aliffatulmf/medotcom/parser"
)

func NewRequest(method, url string, body io.Reader, cookies []parser.Cookie) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Host", "you.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.0.0")

	for _, c := range cookies {
		cookie := c.Cookie()
		if err = cookie.Valid(); err != nil {
			return nil, fmt.Errorf("cookie not valid: %s", err)
		}

		req.Header.Add("Cookie", cookie.String())
	}

	return req, err
}
