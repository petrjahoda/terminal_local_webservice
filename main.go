package main

import (
	"github.com/fatih/color"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"net/http"
	"time"
)

var Interface = "en7"

func main() {
	router := httprouter.New()
	streamer := sse.New()
	router.GET("/", Homepage)
	router.GET("/darcula.css", darcula)
	router.GET("/metro.min.js", metrojs)
	router.GET("/metro-all.css", metrocss)
	router.Handler("GET", "/listen", streamer)
	go StreamTime(streamer)
	_ = http.ListenAndServe(":80", router)
}

func StreamTime(streamer *sse.Streamer) {
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func darcula(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "darcula.css")
}

func metrojs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "metro.min.js")
}
func metrocss(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "metro-all.css")
}

func LogInfo(reference, data string) {
	color.Green(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] " + data)
}

func LogError(reference, data string) {
	color.Red(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] " + data)
}

func LogWarning(reference, data string) {
	color.Yellow(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] " + data)
}
