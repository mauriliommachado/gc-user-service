package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauriliommachado/go-commerce/user-service/data"
	"github.com/mauriliommachado/go-commerce/user-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var pr data.UserRepository

// InitServer function
func InitServer(app *models.App) error {
	log.Println("Initialing server")
	pr.C = app.Collection
	router := httprouter.New()
	router.GET("/ping", ping)
	router.POST("/authenticate", authenticate)
	router.GET("/validateToken/:jwt", validateToken)
	router.GET("/user/:id", get)
	router.GET("/user", getAll)
	router.PUT("/user/:id", protectMiddleware(update))
	router.POST("/user", protectMiddleware(add))
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

func get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p, err := pr.Get(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, _ := json.Marshal(p)
	w.Write(response)
}

func delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := pr.Delete(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func getAll(w http.ResponseWriter, r *http.Request, s httprouter.Params) {
	response, _ := json.Marshal(pr.GetAll())
	w.Write(response)
}

func update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID, _ = primitive.ObjectIDFromHex(ps.ByName("id"))
	pr.Update(&user)
}

func add(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	pr.Create(&user)
	w.Header().Add("Location", user.ID.Hex())
	w.WriteHeader(http.StatusCreated)
}
