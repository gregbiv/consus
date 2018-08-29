package docs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type mockFileSystem struct{}

func (f *mockFileSystem) Open(name string) (http.File, error) {
	dirFoo := &mockFile{
		name: "foo",
		body: []byte("foobody"),
	}
	switch name {
	case "/dir/foo":
		return dirFoo, nil
	case "/dir/":
		return &mockFile{
			name: "/dir/",
			files: []os.FileInfo{
				dirFoo,
				&mockFile{
					name:  "bar",
					files: []os.FileInfo{&mockFile{}},
				},
			},
		}, nil
	}
	return nil, os.ErrNotExist
}

type mockFile struct {
	name   string
	body   []byte
	files  []os.FileInfo
	offset int64
}

// http.File
func (f *mockFile) Readdir(count int) ([]os.FileInfo, error) { return f.files, nil }
func (f *mockFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		f.offset = offset
	case 1:
		f.offset += offset
	case 2:
		f.offset = int64(len(f.body)) - offset
	}
	return f.offset, nil
}
func (f *mockFile) Stat() (os.FileInfo, error) { return f, nil }
func (f *mockFile) Close() error               { return nil }
func (f *mockFile) Read(p []byte) (n int, err error) {
	return bytes.NewBuffer(f.body[f.offset:]).Read(p)
}

// os.FileInfo
func (f *mockFile) Name() string       { return f.name }
func (f *mockFile) Size() int64        { return int64(len(f.body)) }
func (f *mockFile) Mode() os.FileMode  { return 0 }
func (f *mockFile) ModTime() time.Time { return time.Now() }
func (f *mockFile) IsDir() bool        { return len(f.files) > 0 }
func (f *mockFile) Sys() interface{}   { return nil }

func TestDocServer(t *testing.T) {
	t.Parallel()

	router := chi.NewMux()
	DocServer(router, "/", &mockFileSystem{})

	t.Run("Getting a documentation file", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/dir/foo", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Getting a non-existent documentation file", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/dir/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDocServer_WrongConfiguration(t *testing.T) {
	router := chi.NewMux()
	assert.Panics(t, func() {
		DocServer(router, "/{}", &mockFileSystem{})
	})
}
