package static

import (
	"net/http"
	"path"
	"path/filepath"
)

// neutered is a http file system wrapper that disables FileServer Directory Listings
// and roots every path in /static.
type Neutered struct {
	Prefix     string
	FileSystem http.FileSystem
}

func (n Neutered) Open(name string) (http.File, error) {
	name = path.Join(n.Prefix, name)
	f, err := n.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(name, "index.html")
		if _, err := n.FileSystem.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}
	return f, nil
}
