package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/mauriliommachado/go-commerce/user-service/models"
)

func protectMiddleware(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := parseBearerToken(bearerToken[1])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(models.Exception{Message: err.Error()})
					return
				}
				if token.Valid {
					next(w, req, ps)
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(models.Exception{Message: "Invalid Authorization token"})
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Exception{Message: "An Authorization header is required"})
	})
}

func parseBearerToken(bearerToken string) (*jwt.Token, error) {
	return jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("superSecretKey"), nil
	})
}

func signedTokenString(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Unix() + 3600,
	})
	tokenString, err := token.SignedString([]byte("superSecretKey"))
	if err != nil {
		log.Println(err)
	}
	return tokenString
}

func validateCredentials(user models.User) bool {
	if user.Username == user.Password {
		return true
	}
	return false
}
