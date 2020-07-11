package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const basePath = "/api"

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type handler struct {
}

func New() http.Handler {
	h := handler{}
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("%v/users/login", basePath), responseHandler(h.login)).Methods(http.MethodPost)
	return router
}

func responseHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status != http.StatusNoContent {
			if err := json.NewEncoder(w).Encode(response{Data: data, Success: err == nil}); err != nil {
				log.Printf("could not encode response to output: %v", err)
			}
		}
	}
}

func (h *handler) login(w io.Writer, r *http.Request) (interface{}, int, error) {
	var fields map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	username, ok := fields["username"].(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("missing username field")
	}

	password, ok := fields["password"].(string)
	if !ok {
		return nil, http.StatusBadRequest, errors.New("missing password field")
	}

	scopes, err := getMockedData(username, password)
	data := map[string][]string{
		"scopes": scopes,
	}

	return data, http.StatusOK, nil
}

func getMockedData(username, password string) ([]string, error) {
	mockedPassword := "123456789"
	if password != mockedPassword {
		return nil, errors.New("either username or password are wrong")
	}

	switch username {
	case "all_scopes_user":
		return []string{"inventory", "payment", "orders", "other"}, nil
	case "no_scopes_user":
		return []string{}, nil
	case "inventory_scopes_user":
		return []string{"inventory"}, nil
	default:
		return nil, errors.New("either username or password are wrong")
	}
}
