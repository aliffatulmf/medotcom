package parser

import (
	"os"
	"testing"
)

func createTempFile(content string, t *testing.T) string {
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}

	return tempFile.Name()
}

func TestFileParser_Json(t *testing.T) {
	validJson := createTempFile(`[{"Domain":"example.com","HostOnly":true,"Path":"/","Secure":true,"Name":"test","Value":"value"}]`, t)
	invalidJson := createTempFile(`invalid json`, t)
	invalidDataJson := createTempFile(`{"Domain":"example.com","HostOnly":true,"Path":"/","Secure":true,"Name":"test"}`, t)

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "File exists and contains valid JSON",
			file:    validJson,
			wantErr: false,
		},
		{
			name:    "File exists but contains invalid JSON",
			file:    invalidJson,
			wantErr: true,
		},
		{
			name:    "File does not exist",
			file:    "nonexistent.json",
			wantErr: true,
		},
		{
			name:    "File exists and contains invalid data",
			file:    invalidDataJson,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewParser(tt.file)
			_, err := f.Json()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileParser.Json() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileParser_Text(t *testing.T) {
	validText := createTempFile("example.com\tTRUE\t/\tTRUE\t0\ttest\tvalue", t)
	invalidText := createTempFile("invalid text", t)
	invalidDataText := createTempFile("example.com\tTRUE\t/\tTRUE\t0\ttest", t)

	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "File exists and contains valid Text",
			file:    validText,
			wantErr: false,
		},
		{
			name:    "File exists but contains invalid Text",
			file:    invalidText,
			wantErr: true,
		},
		{
			name:    "File does not exist",
			file:    "nonexistent.txt",
			wantErr: true,
		},
		{
			name:    "File exists and contains invalid data",
			file:    invalidDataText,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewParser(tt.file)
			_, err := f.Text()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileParser.Text() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
