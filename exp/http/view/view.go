package view

import "net/http"

type Renderer interface {
	http.Handler
	Render(http.ResponseWriter, *http.Request, any)
}
