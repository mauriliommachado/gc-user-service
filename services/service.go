package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauriliommachado/go-commerce/user-service/models"
)

// InitServer function
func InitServer(app *models.App) error {
	log.Println("Initialing server")
	router := httprouter.New()
	router.GET("/ping", ping)
	router.GET("/protected", protectMiddleware(http.HandlerFunc(protected)))
	router.POST("/authenticate", authenticate)
	router.GET("/validateToken/:jwt", validateToken)
	log.Fatal(http.ListenAndServe(":8080", router))
	app.Router = router
	return nil
}

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "pong\n")
}

func protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You have the key!\n")
}

func authenticate(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(models.Exception{Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if validateCredentials(user) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.JwtToken{Token: signedTokenString(models.User{Username: user.Username})})
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(models.Exception{Message: "invalid credentials"})
}

func validateToken(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	token, err := parseBearerToken(ps.ByName("jwt"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Exception{Message: "invalid token"})
		return
	}
	if token.Valid {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(models.Exception{Message: "invalid token"})
}
