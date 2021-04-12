package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os/exec"
	"time"
)

func Screenshot(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Screenshot loading")
	start := time.Now()
	command := "sudo"
	args := []string{"-u", "pi", "maim", "image.png"}
	argumentDebug := ""
	for _, arg := range args {
		argumentDebug += arg + " "
	}
	LogInfo("MAIN", command+" "+argumentDebug)
	data, err := exec.Command(command, args...).Output()
	if err != nil {
		LogError("MAIN", err.Error())
	}
	LogInfo("MAIN", "Screenshot taken: "+string(data))
	HomepageLoaded = false
	renderTemplate(w, "screenshot", &Page{})
	LogInfo("MAIN", "Screenshot loaded in "+time.Since(start).String())
}
