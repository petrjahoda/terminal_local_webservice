package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kbinani/screenshot"
	"image/png"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func Screenshot(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Screenshot loading")
	start := time.Now()
	command := "sudo"
	args := []string{"-u", "zapsi", "maim", "image.png"}
	if runtime.GOOS == "linux" {
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
	} else {
		n := screenshot.NumActiveDisplays()
		for i := 0; i < n; i++ {
			img, err := screenshot.CaptureDisplay(i)
			if err != nil {
				LogError("MAIN", "Error generating screenshot: "+err.Error())
				continue
			}
			fileName := "image.png"
			file, _ := os.Create(fileName)
			_ = png.Encode(file, img)
			LogInfo("MAIN", "Generated screenshot: "+fileName)
			file.Close()
		}

	}
	HomepageLoaded = false
	renderTemplate(w, "screenshot", &Page{})
	LogInfo("MAIN", "Screenshot loaded in "+time.Since(start).String())
}
