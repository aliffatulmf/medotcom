package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aliffatulmf/medotcom/util"
)

const (
	TEXT = ".txt"
	JSON = ".json"
)

type FileParser struct {
	File string
}

func NewParser(f string) ([]Cookie, error) {
	if err := util.IsFileCorrect(f); err != nil {
		return nil, fmt.Errorf("error checking file: %w", err)
	}

	fp := FileParser{File: f}

	switch fp.Ext() {
	case JSON:
		return fp.Json()
	case TEXT:
		return fp.Text()
	default:
		return nil, fmt.Errorf("unsupported file extension for file %s", f)
	}
}

func (f *FileParser) Ext() string {
	return filepath.Ext(f.File)
}

func (f *FileParser) Json() ([]Cookie, error) {
	var cookies []Cookie

	byteFile, err := os.ReadFile(f.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", f.File, err)
	}

	if err = json.Unmarshal(byteFile, &cookies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from file %s: %w", f.File, err)
	}

	return cookies, nil
}

func (f *FileParser) Text() ([]Cookie, error) {
	file, err := os.Open(f.File)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", f.File, err)
	}
	defer file.Close()

	var cookies []Cookie
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// skip if line is empty or starts with '#'
		if line == "" || line[0] == '#' {
			continue
		}

		cookie, err := parseLineArray(strings.Fields(line))
		if err != nil {
			return nil, fmt.Errorf("failed to parse line in file %s: %w", f.File, err)
		}

		cookies = append(cookies, *cookie)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while scanning file %s: %w", f.File, err)
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

	expire, err := strconv.ParseFloat(s[4], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing expiration date: %w", err)
	}

	return &Cookie{
		Domain:         s[0],
		HostOnly:       hostOnly,
		Path:           s[2],
		Secure:         secure,
		ExpirationDate: expire,
		Name:           s[5],
		Value:          s[6],
	}, nil
}
