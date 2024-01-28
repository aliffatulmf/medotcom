package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aliffatulmf/medotcom/parser"
	"github.com/aliffatulmf/medotcom/requests"
)

const (
	GET    = "GET"
	DELETE = "DELETE"
)

var (
	cookieFile   string
	actionMethod string
)

func init() {
	parser := flag.NewFlagSet("You Connect", flag.ExitOnError)
	parser.StringVar(&cookieFile, "cookie", "", "Path to the cookie file. This file should contain the cookies required for authentication.")
	parser.StringVar(&actionMethod, "action", GET, "Action to perform. Can be either 'GET' or 'DELETE'.")
	if err := parser.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	cookies, err := parser.NewParser(cookieFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res := getChat(cookies)

	switch strings.ToUpper(actionMethod) {
	case GET:
		showChatTitle(res)
	case DELETE:
		deleteAllChat(res, cookies)
	default:
		fmt.Println("Invalid action method.")
		os.Exit(1)
	}
}

func getChat(cookies []parser.Cookie) *requests.ChatResponse {
	result, err := requests.RequestGET(&requests.RequestOptions{
		Payload: strings.NewReader(`{"count": 0}`),
		Cookies: cookies,
	})

	if err != nil {
		fmt.Printf("Failed to fetch chat from database: %s\n", err)
		os.Exit(1)
	}

	if result == nil || len(result.Chats) == 0 {
		fmt.Println("No chat available.")
		os.Exit(0)
	}

	return result
}

func showChatTitle(cr *requests.ChatResponse) {
	for _, chat := range cr.Chats {
		fmt.Println(chat.Title)
	}
}

func deleteAllChat(cr *requests.ChatResponse, cookies []parser.Cookie) {
	n := 5
	if len(cr.Chats) < 5 {
		n = len(cr.Chats)
	}

	lim := make(chan struct{}, n)

	for _, chat := range cr.Chats {
		lim <- struct{}{}

		go func(chat requests.Chat) {
			body := fmt.Sprintf(`{"chatId": "%s"}`, chat.ChatID)

			err := requests.RequestDELETE(&requests.RequestOptions{
				Payload: strings.NewReader(body),
				Cookies: cookies,
			})

			if err == nil {
				fmt.Println("Successfully deleted chat:", chat.Title)
			} else {
				fmt.Println("Failed to delete chat with title:", chat.Title)
			}

			<-lim
		}(chat)
	}
}
