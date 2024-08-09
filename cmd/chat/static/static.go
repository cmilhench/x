package static

import (
	"embed"
	"net/http"
	"path"
	"path/filepath"
)

//go:embed files
var static embed.FS

func Handler() http.Handler {
	return http.FileServer(neutered{http.FS(static)})
}

// neutered is a http file system wrapper that disables FileServer Directory Listings
// and roots every path in /static
type neutered struct {
	fs http.FileSystem
}

func (n neutered) Open(name string) (http.File, error) {
	name = path.Join("files", name)
	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}
	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(name, "index.html")
		if _, err := n.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}
	return f, nil
}
