package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os/exec"
)

func screenshotPage(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	command := "sudo"
	args := []string{"-u", "pi", "maim", "image.png"}
	argumentDebug := ""
	for _, arg := range args {
		argumentDebug += arg + " "
	}
	exec.Command(command, args...)
	renderTemplate(w, "screenshot", &Page{})
}
