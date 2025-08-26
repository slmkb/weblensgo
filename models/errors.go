package models

import (
	"fmt"
	"io"
	"net/http"
)

type FIleError struct {
	Issue string
}

func (fe FIleError) Error() string {
	return fmt.Sprintf("invalid file type: %v", fe.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("check content type: %w", err)
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("check content type: %w", err)
	}
	contentType := http.DetectContentType(testBytes)
	for _, t := range allowedTypes {
		if t == contentType {
			return nil
		}
	}
	return FIleError{
		Issue: fmt.Sprintf("check content type: %v", contentType),
	}
}
