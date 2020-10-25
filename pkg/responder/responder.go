package responder

import (
	"encoding/json"
	"log"
	"net/http"
)

var errFailedEncoding = []byte(`{"error:"Failed to encode response."}`)

type responder struct {
	http.ResponseWriter
	req *http.Request
}

type Responder interface {
	OK(message string, data interface{})
	BadRequest(message string, data interface{})
}

type jsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func New(w http.ResponseWriter, req *http.Request) Responder {
	return &responder{
		ResponseWriter: w,
		req:            req,
	}
}

func (w *responder) writeJSON(body interface{}) {
	bs, err := json.Marshal(body)
	if err != nil {
		log.Println("Something went wrong when encoding JSON output: " + err.Error())
		bs = errFailedEncoding
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = w.Write(bs)
	if err != nil {
		log.Println("Failed to write HTTP response: " + err.Error())
	}
}

func (w *responder) writeText(message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(message)); err != nil {
		log.Println("Failed to write HTTP response: " + err.Error())
	}
}

func (w *responder) write(status int, message string, data interface{}) {
	header := w.req.Header.Get("accept")
	if header == "" {
		header = "text/plain"
	}

	w.WriteHeader(status)

	accept := http.DetectContentType([]byte(header))
	switch accept {
	case "application/json":
		w.writeJSON(&jsonResponse{
			Message: message,
			Data:    data,
		})

	default:
		w.writeText(message)
	}
}

func (w *responder) OK(message string, data interface{}) {
	w.write(http.StatusOK, message, data)
}

func (w *responder) BadRequest(message string, data interface{}) {
	w.write(http.StatusBadRequest, message, data)
}
