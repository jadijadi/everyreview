package media

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Minimal valid PNG header so DetectContentType reports image/png.
var pngHeader = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 13, 'I', 'H', 'D', 'R'}

func TestSaveProductImage_WritesFileUnderRandomKey(t *testing.T) {
	store, err := NewDiskStore(t.TempDir())
	require.NoError(t, err)
	payload := append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{1}, 1000)...)

	key, err := store.SaveProductImage(context.Background(), bytes.NewReader(payload))

	require.NoError(t, err)
	require.Regexp(t, keyPattern, key)
	require.True(t, filepath.IsLocal(key))
	written, err := os.ReadFile(store.path(key))
	require.NoError(t, err)
	require.Equal(t, payload, written)

	require.NoError(t, store.Delete(context.Background(), key))
	_, err = os.Stat(store.path(key))
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestSaveProductImage_RejectsNonImages(t *testing.T) {
	store, err := NewDiskStore(t.TempDir())
	require.NoError(t, err)

	_, err = store.SaveProductImage(context.Background(), bytes.NewReader([]byte("<html>not an image</html>")))

	require.ErrorIs(t, err, ErrUnsupportedType)
}

func TestSaveProductImage_RejectsOversizedAndCleansUp(t *testing.T) {
	dir := t.TempDir()
	store, err := NewDiskStore(dir)
	require.NoError(t, err)
	payload := append(append([]byte{}, pngHeader...), bytes.Repeat([]byte{1}, MaxImageBytes)...)

	_, err = store.SaveProductImage(context.Background(), bytes.NewReader(payload))

	require.ErrorIs(t, err, ErrTooLarge)
	entries, err := os.ReadDir(filepath.Join(dir, "products"))
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestDelete_RefusesKeysOutsideTheStore(t *testing.T) {
	store, err := NewDiskStore(t.TempDir())
	require.NoError(t, err)

	require.Error(t, store.Delete(context.Background(), "../etc/passwd"))
}
