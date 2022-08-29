package util

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Routes map[string]http.HandlerFunc
type RequestValidator func(http.ResponseWriter, *http.Request) error

func composeValidator(wrappers []RequestValidator) RequestValidator {
	var wfunc RequestValidator
	for _, wf := range wrappers {
		if wfunc == nil {
			wfunc = wf
			continue
		}
		f := func(w http.ResponseWriter, r *http.Request) error {
			if err := wfunc(w, r); err != nil {
				return err
			}
			return wf(w, r)
		}
		wfunc = f
	}
	return wfunc
}

func BuildRoutes(r *mux.Router, routes Routes, wrappers ...RequestValidator) {
	validator := composeValidator(wrappers)
	for k, f := range routes {
		r.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			if err := validator(w, r); err != nil {
				return
			}
			f(w, r)
		})
	}
}
