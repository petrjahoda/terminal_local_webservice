package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"time"
)

var Interface = "Ethernet"

const deleteLogsAfter = 240 * time.Hour

type Page struct {
	Title string
	Body  []byte
}

func main() {
	LogDirectoryFileCheck("MAIN")
	CreateConfigIfNotExists()
	router := httprouter.New()
	streamer := sse.New()
	router.GET("/", Homepage)
	router.GET("/screenshot", Screenshot)
	router.GET("/changenetwork", ChangeNetwork)
	router.GET("/changenetworktodhcp", ChangeNetworkToDhcp)
	router.GET("/restart", Restart)
	router.GET("/shutdown", Shutdown)
	router.GET("/setup", Setup)
	router.GET("/css/darcula.css", darcula)
	router.GET("/js/metro.min.js", metrojs)
	router.GET("/css/metro-all.css", metrocss)
	router.GET("/image.png", image)
	router.Handler("GET", "/listen", streamer)
	go StreamTime(streamer)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":8000", router)
}

func StreamTime(streamer *sse.Streamer) {
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("html/" + tmpl + ".html")
	_ = t.Execute(w, p)
}

func darcula(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "css/darcula.css")
}

func metrojs(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "js/metro.min.js")
}

func metrocss(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "css/metro-all.css")
}

func image(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "image.png")
}
