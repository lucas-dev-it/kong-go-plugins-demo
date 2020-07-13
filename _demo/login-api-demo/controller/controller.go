package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	basePath = "/api"
	signKey  = "someFakeTokenJustForDemo"
)

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
	router.HandleFunc(fmt.Sprintf("%v/users/test-kong", basePath), responseHandler(h.testKongJWTPlugin)).Methods(http.MethodGet)
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

func (h *handler) testKongJWTPlugin(w io.Writer, r *http.Request) (interface{}, int, error) {
	return "if you see this message it's because of kong pass the request through", 200, nil
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
	token, err := buildJWT(scopes)

	data := map[string]string{
		"accessToken": token,
	}

	return data, http.StatusOK, nil
}

func buildJWT(scopes []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{}
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(time.Minute * 15).Unix()
	claims["scopes"] = scopes
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func getMockedData(username, password string) ([]string, error) {
	mockedPassword := "123456789"
	if password != mockedPassword {
		return nil, errors.New("either username or password are wrong")
	}

	switch username {
	case "all_scopes_user":
		return []string{"inventory", "payment", "order", "other"}, nil
	case "no_scopes_user":
		return []string{}, nil
	case "inventory_scopes_user":
		return []string{"inventory"}, nil
	default:
		return nil, errors.New("either username or password are wrong")
	}
}
