package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type File struct {
	ID       string
	Content  []byte
	Hash     []byte
	Bucket   string
	MimeType string
	Name     string
}

func (f *File) Stem() string {
	var ext string
	parts := strings.Split(f.Name, ".")
	if len(parts) > 1 {
		ext = parts[1]
	}
	return fmt.Sprintf("%s.%s", f.ID, ext)
}

// Size returns file total size in bytes.
func (f *File) Size() uint64 {
	return uint64(len(f.Content))
}

func (f *File) SaveTo(ctx context.Context, s Storage, path string) error {
	uploadCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	if len(f.Content) == 0 {
		return errors.New("file content is empty")
	}

	mimeType := http.DetectContentType(f.Content)
	allowedTypes := map[string]bool{
		"image/jpeg":        true,
		"image/png":         true,
		"image/gif":         true,
		"application/pdf":   true,
		"application/dicom": true,
	}

	if !allowedTypes[mimeType] {
		return fmt.Errorf("invalid file type: %s", mimeType)
	}

	contentReader := bytes.NewReader(f.Content)
	key := fmt.Sprintf("%s/%s", path, f.Stem())
	bucket, err := s.Upload(uploadCtx, key, contentReader, nil)
	if err != nil {
		return err
	}

	f.Bucket = bucket
	f.MimeType = mimeType
	return nil
}
