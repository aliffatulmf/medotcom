package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	TEXT = ".txt"
	JSON = ".json"
)

type FileParser struct {
	path string
}

func NewParser(f string) *FileParser {
	return &FileParser{path: f}
}

func (f *FileParser) Validate() error {
	if f.path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if strings.Contains(f.path, "..") {
		return fmt.Errorf("path traversal detected")
	}

	if _, err := os.Stat(f.path); err != nil {
		return fmt.Errorf("file %s not found: %w", f.path, err)
	}

	return nil
}

func (f *FileParser) Parse() ([]Cookie, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}

	switch filepath.Ext(f.path) {
	case JSON:
		return f.Json()
	case TEXT:
		return f.Text()
	default:
		return nil, fmt.Errorf("unsupported file extension")
	}
}

func (f *FileParser) readFile() (io.ReadCloser, error) {
	fopen, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}

	stat, err := fopen.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() == 0 || stat.Size() > 15e3 { // 15kb
		return nil, fmt.Errorf("file size is too large or empty")
	}

	return fopen, nil
}

func (f *FileParser) Json() ([]Cookie, error) {
	fopen, err := f.readFile()
	if err != nil {
		return nil, fmt.Errorf("readFile: %s", err)
	}
	defer fopen.Close()

	var cookies []Cookie
	if err := json.NewDecoder(fopen).Decode(&cookies); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return cookies, nil
}

func (f *FileParser) Text() ([]Cookie, error) {
	fopen, err := f.readFile()
	if err != nil {
		return nil, fmt.Errorf("fileOpen: %s", err)
	}
	defer fopen.Close()

	scanner := bufio.NewScanner(fopen)

	var cookies []Cookie
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || line[0] == '#' {
			continue
		}

		cookie, err := parseLineArray(strings.Fields(line))
		if err != nil {
			return nil, fmt.Errorf("failed to parse line in file %s: %w", f.path, err)
		}

		cookies = append(cookies, *cookie)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while scanning file %s: %w", f.path, err)
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
