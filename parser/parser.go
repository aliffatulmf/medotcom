package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileParser struct {
	Path string
}

func NewParser(cookie string) (*FileParser, error) {
	cookie = strings.TrimSpace(cookie)

	if cookie == "" {
		return nil, fmt.Errorf("NewParser: cookie file path cannot be empty")
	}

	stat, err := os.Lstat(cookie)
	if err != nil {
		return nil, fmt.Errorf("NewParser: %w", err)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("NewParser: cookie file path is a directory")
	}

	if stat.Size() == 0 {
		return nil, fmt.Errorf("NewParser: cookie file is empty")
	}

	if !strings.HasSuffix(cookie, TEXT) && !strings.HasSuffix(cookie, JSON) {
		return nil, fmt.Errorf("NewParser: unsupported file extension")
	}

	return &FileParser{Path: cookie}, nil
}

func (f *FileParser) Parse() ([]Cookie, error) {
	switch filepath.Ext(f.Path) {
	case JSON:
		return f.Json()
	case TEXT:
		return f.Text()
	default:
		return nil, fmt.Errorf("unsupported file extension")
	}
}

func (f *FileParser) Json() ([]Cookie, error) {
	fopen, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer fopen.Close()

	if err = CheckFileSize(fopen, 15*1024); err != nil {
		return nil, fmt.Errorf("file size limit: %w", err)

	}

	var cookies []Cookie
	if err := json.NewDecoder(fopen).Decode(&cookies); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return cookies, nil
}

func (f *FileParser) Text() ([]Cookie, error) {
	fopen, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer fopen.Close()

	if err = CheckFileSize(fopen, 15*1024); err != nil {
		return nil, fmt.Errorf("file size limit: %w", err)
	}

	scanner := bufio.NewScanner(fopen)

	var cookies []Cookie
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || line[0] == '#' {
			continue
		}

		cookie, err := parseLineArray(strings.Fields(line))
		if err != nil {
			return nil, err
		}

		cookies = append(cookies, *cookie)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while scanning file %s: %w", f.Path, err)
	}

	return cookies, nil
}

func parseLineArray(s []string) (*Cookie, error) {
	if len(s) < 7 {
		return nil, fmt.Errorf("insufficient data: expected 7 fields, got %d", len(s))
	}

	hostOnly, err := strconv.ParseBool(s[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing hostonly: %w", err)
	}

	secure, err := strconv.ParseBool(s[3])
	if err != nil {
		return nil, fmt.Errorf("error parsing secure: %w", err)
	}

	return &Cookie{
		Domain:   s[0],
		HostOnly: hostOnly,
		Path:     s[2],
		Secure:   secure,
		Name:     s[5],
		Value:    s[6],
	}, nil
}
