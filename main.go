package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const version = "2021.2.1.19"
const programName = "Terminal local webservice"
const programDesription = "Display local web for rpi terminals"

type Page struct {
	Title string
	Body  []byte
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	router := httprouter.New()
	networkDataStreamer := sse.New()
	router.GET("/image.png", image)
	router.GET("/", indexPage)
	router.GET("/screenshot", screenshotPage)
	router.GET("/setup", setupPage)
	router.POST("/password", checkPassword)
	router.POST("/restart", restartRpi)
	router.POST("/shutdown", shutdownRpi)
	router.POST("/dhcp", changeToDhcp)
	router.POST("/static", changeToStatic)
	router.ServeFiles("/font/*filepath", http.Dir("font"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.Handler("GET", "/networkdata", networkDataStreamer)
	go StreamNetworkData(networkDataStreamer)
	_ = http.ListenAndServe(":9999", router)
}

func LoadSettingsFromConfigFile() string {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	readFile, _ := ioutil.ReadFile(configFullPath)
	ConfigFile := ServerIpAddress{}
	_ = json.Unmarshal(readFile, &ConfigFile)
	ServerIpAddress := ConfigFile.ServerIpAddress
	return ServerIpAddress
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	serviceConfig := &service.Config{
		Name:        programName,
		DisplayName: programName,
		Description: programDesription,
	}
	prg := &program{}
	s, _ := service.New(prg, serviceConfig)
	_ = s.Run()
}

func StreamNetworkData(streamer *sse.Streamer) {
	timeToSend := "20"
	for {
		interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
		interfaceServerIpAddress := LoadSettingsFromConfigFile()
		serverAccessible := CheckServerIpAddress(interfaceServerIpAddress)
		if !serverAccessible {
			interfaceServerIpAddress = interfaceServerIpAddress + ", offline"
		}
		streamer.SendString("", "networkdata", interfaceIpAddress+";"+interfaceMask+";"+interfaceGateway+";"+dhcpEnabled+";"+timeToSend+";"+interfaceServerIpAddress+";"+interfaceServerIpAddress)
		time.Sleep(5 * time.Second)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("html/" + tmpl + ".html")
	_ = t.Execute(w, p)
}
func image(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "image.png")
}
