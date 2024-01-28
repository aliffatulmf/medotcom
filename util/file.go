package util

import (
	"fmt"
	"os"
	"strings"
)

func IsFileCorrect(f string) error {
	if len(strings.TrimSpace(f)) == 0 {
		return fmt.Errorf("error path is not specified")
	}

	fi, err := os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", err)
		}

		if pe, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("failed to read file: %s", pe.Err)
		}
	}

	if fi.IsDir() {
		return fmt.Errorf("error file is dir")
	}

	return nil
}
