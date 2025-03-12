package utils

import "os"

func OpenCreate(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		Fatal("Error opening file:", err)
		return nil
	}

	return file
}
