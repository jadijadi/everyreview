// Package media stores uploaded images. MVP simplification (see backend/README.md
// "MVP scope"): files live on local disk and are served by the API itself, instead
// of the S3 presigned-upload + CDN flow from ADR-0011 — there is no AWS yet.
package media

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const MaxImageBytes = 5 << 20

var (
	ErrTooLarge        = errors.New("media: image must be at most 5 MB")
	ErrUnsupportedType = errors.New("media: image must be JPEG, PNG or WebP")
)

var extensionByContentType = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

// keyPattern is the only shape of key this store ever produces, and the only one
// it will serve — it doubles as the path-traversal guard.
var keyPattern = regexp.MustCompile(`^products/[0-9a-f]{32}\.(jpg|png|webp)$`)

type Store interface {
	SaveProductImage(ctx context.Context, r io.Reader) (key string, err error)
	Delete(ctx context.Context, key string) error
}

type DiskStore struct {
	dir string
}

func NewDiskStore(dir string) (*DiskStore, error) {
	if err := os.MkdirAll(filepath.Join(dir, "products"), 0o755); err != nil {
		return nil, fmt.Errorf("media: create dir: %w", err)
	}
	return &DiskStore{dir: dir}, nil
}

// SaveProductImage sniffs the real content type (the client's declared type is
// untrusted), enforces the size cap and writes the file under a random key.
func (s *DiskStore) SaveProductImage(ctx context.Context, r io.Reader) (string, error) {
	head := make([]byte, 512)
	n, err := io.ReadFull(r, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("media: read: %w", err)
	}
	head = head[:n]
	ext, ok := extensionByContentType[http.DetectContentType(head)]
	if !ok {
		return "", ErrUnsupportedType
	}

	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("media: random key: %w", err)
	}
	key := fmt.Sprintf("products/%s.%s", hex.EncodeToString(raw[:]), ext)

	f, err := os.OpenFile(s.path(key), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("media: create file: %w", err)
	}
	// +1 so we can tell "exactly the limit" from "over the limit".
	limited := io.LimitReader(io.MultiReader(bytes.NewReader(head), r), MaxImageBytes+1)
	written, copyErr := io.Copy(f, limited)
	closeErr := f.Close()
	switch {
	case copyErr != nil:
		_ = os.Remove(s.path(key))
		return "", fmt.Errorf("media: write: %w", copyErr)
	case closeErr != nil:
		_ = os.Remove(s.path(key))
		return "", fmt.Errorf("media: close: %w", closeErr)
	case written > MaxImageBytes:
		_ = os.Remove(s.path(key))
		return "", ErrTooLarge
	}
	return key, nil
}

func (s *DiskStore) Delete(ctx context.Context, key string) error {
	if !keyPattern.MatchString(key) {
		return fmt.Errorf("media: refusing to delete invalid key %q", key)
	}
	if err := os.Remove(s.path(key)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("media: delete: %w", err)
	}
	return nil
}

func (s *DiskStore) path(key string) string {
	return filepath.Join(s.dir, filepath.FromSlash(key))
}
