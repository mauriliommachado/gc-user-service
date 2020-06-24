package main

import (
	"log"

	"github.com/mauriliommachado/go-commerce/user-service/models"
	"github.com/mauriliommachado/go-commerce/user-service/services"
)

func main() {
	log.Println("Initialing application")
	// Initialize server
	app := models.App{}
	log.Println("Loading configurations")
	services.InitServer(&app)
}
