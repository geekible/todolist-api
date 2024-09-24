package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"geekible.todolist/src/domain"
)

type JsonHelpers struct{}

func InitJsonHelpers() *JsonHelpers {
	return &JsonHelpers{}
}

func (h *JsonHelpers) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (h *JsonHelpers) WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

func (h *JsonHelpers) ErrorJSON(w http.ResponseWriter, err error, status int, errSource string) error {
	statusCode := http.StatusBadRequest

	if status > 0 {
		statusCode = status
	}

	var payload domain.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return h.WriteJSON(w, statusCode, payload)
}
