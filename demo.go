package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

type DemoData struct {
	Version string
}

func demo1(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_1.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo2(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_2.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo3(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_3.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo4(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_4.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo5(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_5.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo6(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_6.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo7(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_7.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo8(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_8.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo9(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_9.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}

func demo10(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	tmpl := template.Must(template.ParseFiles("html/demo_10.html"))
	data := DemoData{
		Version: version,
	}
	_ = tmpl.Execute(writer, data)
}
