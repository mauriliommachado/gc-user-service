package models

import "github.com/julienschmidt/httprouter"

// App application manager
type App struct {
	Router *httprouter.Router
}
