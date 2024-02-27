package requests

import (
	"errors"
)

var ErrNoChatFound = errors.New("No chat found.")

type Chat struct {
	ID         string `json:"_id"`
	ChatID     string `json:"chat_id"`
	Title      string `json:"title"`
	DateUpdate string `json:"date_updated"`
}

type ChatResponse struct {
	Chats []Chat `json:"chats"`
}

func (c *ChatResponse) Len() int {
	if c == nil || c.Chats == nil {
		return 0
	}
	return len(c.Chats)
}

func (c *ChatResponse) Empty() bool {
	return c.Len() == 0
}
