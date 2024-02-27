package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aliffatulmf/medotcom/parser"
	"github.com/aliffatulmf/medotcom/requests"
)

var (
	cookieFlag string
	deleteFlag bool
)

func init() {
	flagSet := flag.NewFlagSet("You Connect", flag.ExitOnError)
	flagSet.StringVar(&cookieFlag, "cookie", "", "path to the cookie file.")
	flagSet.BoolVar(&deleteFlag, "delete", false, "delete chat(s).")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Printf("error parsing flags: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	cookie, err := parser.NewParser(cookieFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	parsed, err := cookie.Parse()
	if err != nil {
		fmt.Printf("failed to parse cookie file: %s\n", err)
		return
	}

	resp, err := requests.RequestGET(parsed)
	if err != nil {
		if errors.Is(err, requests.ErrNoChatFound) {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(err)
		return
	}

	for _, chat := range resp.Chats {
		if deleteFlag {
			if err := requests.RequestDELETE(parsed, &chat); err != nil {
				fmt.Printf("%s failed to delete.\n", chat.Title)
				continue
			}

			fmt.Printf("%s deleted.\n", chat.Title)
		} else {
			fmt.Printf("Title: %s\n", chat.Title)
		}
	}
}
