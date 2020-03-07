package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func Password(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Password loading")
	start := time.Now()
	HomepageLoaded = false
	renderTemplate(w, "password", &Page{})
	LogInfo("MAIN", "Password loaded in "+time.Since(start).String())
}
