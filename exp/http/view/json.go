package view

import (
	"encoding/json"
	"net/http"

	"github.com/cmilhench/x/exp/logger"
)

var _ Renderer = (*jsonView)(nil)

type jsonView struct {
	pretty bool
}

func (v *jsonView) Render(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "application/json")
	var out []byte
	var err error
	if v != nil && v.pretty {
		out, err = json.MarshalIndent(data, "", "  ")
	} else {
		out, err = json.Marshal(data)
	}
	if err != nil {
		logger.Error("failed to marshal json: %+v\n", err)
		http.Error(w, `{"error":"Something went wrong."}`, http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

func (v *jsonView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}
