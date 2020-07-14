package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kong/go-pdk"
	"github.com/dgrijalva/jwt-go"
)

const signKey = "someFakeTokenJustForDemo"

type Conf struct {
	AllowedScopes string `json:"allowed_scopes"`
}

func New() interface{} {
	return &Conf{}
}

func (conf Conf) Access(kong *pdk.PDK) {
	authHeader, err := kong.Request.GetHeader("Authorization")
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(http.StatusUnauthorized, err.Error(), nil)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) < 2 {
		kong.Log.Err("invalid token format")
		kong.Response.Exit(http.StatusBadRequest, "invalid token format", nil)
	}

	if err := conf.verifyToken(parts[1]); err != nil {
		kong.Response.Exit(http.StatusUnauthorized, err.Error(), nil)
	}

	kong.Log.Info("authorized request")
}

func (conf Conf) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return err
	}

	allowedScopesMap := scopesMap(conf.AllowedScopes)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token")
	}

	scopesClaim, ok := claims["scopes"]
	if !ok {
		return errors.New("scopes not present within the JWT token")
	}

	scopesClaimSlice, ok := scopesClaim.([]string)
	if !ok {
		return errors.New("scopes not present within the JWT token")
	}

	for _, sc := range scopesClaimSlice {
		if _, ok := allowedScopesMap[sc]; ok {
			continue
		}
		return fmt.Errorf("%v scope not allowed", sc)
	}

	return nil
}

func scopesMap(allowedScopes string) map[string]struct{} {
	scopes := strings.Split(allowedScopes, ",")

	scopesMap := make(map[string]struct{}, len(scopes))

	for _, s := range scopes {
		scopesMap[s] = struct{}{}
	}

	return scopesMap
}
