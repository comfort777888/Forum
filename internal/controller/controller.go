package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (h *handler) CheckMethod(path, method string, r *http.Request) int {
	if r.Method != method {
		return http.StatusMethodNotAllowed
	}
	if r.URL.Path != path {
		return http.StatusNotFound
	}
	return 200
}

type Error struct {
	ErrorMessage string
	StatusText   string
	Status       int
}

func (h *handler) errorHandler(w http.ResponseWriter, status int, errMessage string) {
	errHandler := h.setError(status, errMessage)
	w.WriteHeader(status)

	log.Printf("error occured: %s", errMessage)

	html, err := template.ParseFiles("ui/template/error.html")
	if err != nil {
		log.Printf("error parse file: %v", err)
		return
	}

	err = html.Execute(w, errHandler)
	if err != nil {
		log.Printf("error execute template:  %v", err)
		return
	}
	// err := h.execute(w, "ui/template/error.html", errHandler)
	// if err != nil {
	// 	fmt.Printf("error execute template: %v", err)
	// 	return
	// }
}

func (h *handler) setError(status int, err string) *Error {
	return &Error{
		ErrorMessage: err,
		Status:       status,
		StatusText:   http.StatusText(status),
	}
}

func (h *handler) execute(w http.ResponseWriter, parse string, data interface{}) error {
	html, err := template.ParseFiles(parse)
	if err != nil {
		log.Printf("error parse file: %v", err)
		return fmt.Errorf("error parse files: %w", err)
	}

	err = html.Execute(w, data)
	if err != nil {
		log.Printf("error execute template: %s %v", parse, err)
		return fmt.Errorf("error execute temp;ate: %w", err)
	}

	return nil
}
