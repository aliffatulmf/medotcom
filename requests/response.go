package requests

type Chat struct {
	ID         string `json:"_id"`
	ChatID     string `json:"chat_id"`
	Title      string `json:"title"`
	DateUpdate string `json:"date_updated"`
}

type ChatResponse struct {
	Chats []Chat `json:"chats"`
}

func (c *ChatResponse) ListID() (result []string) {
	for _, chat := range c.Chats {
		result = append(result, chat.ChatID)
	}
	return
}
