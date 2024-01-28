package parser

import (
	"net/http"
	"strings"
	"time"
)

type Cookie struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate"`
	HostOnly       bool    `json:"hostOnly"`
	HttpOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreId        string  `json:"storeId"`
	Value          string  `json:"value"`
}

func (c Cookie) expireTime() time.Time {
	return time.Unix(int64(c.ExpirationDate), 0).UTC()
}

func (c Cookie) Cookie() *http.Cookie {
	return &http.Cookie{
		Domain:   c.Domain,
		Expires:  c.expireTime(),
		HttpOnly: c.HttpOnly,
		Name:     c.Name,
		Path:     c.Path,
		Secure:   c.Secure,
		Value:    c.Value,
	}
}

func (c Cookie) String() string {
	return strings.Join([]string{c.Name, c.Value}, "=")
}
