package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const version = "2022.1.2.14"
const programName = "Terminal local webservice"
const programDescription = "Display local web for rpi terminals"

var initiated = false
var homepageLoaded = false

type Page struct {
	Title string
	Body  []byte
}

var streamCanRun = false
var streamSync sync.RWMutex

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	serviceConfig := &service.Config{
		Name:        programName,
		DisplayName: programName,
		Description: programDescription,
	}
	prg := &program{}
	s, _ := service.New(prg, serviceConfig)
	_ = s.Run()
}

func (p *program) run() {
	router := httprouter.New()
	networkDataStreamer := sse.New()
	router.GET("/image.png", image)
	router.GET("/", indexPage)
	router.GET("/screenshot", screenshotPage)
	router.GET("/setup", setupPage)
	router.GET("/setup-remote", setupRemotePage)
	router.GET("/demo_1", demo1)
	router.GET("/demo_2", demo2)
	router.GET("/demo_3", demo3)
	router.GET("/demo_4", demo4)
	router.GET("/demo_5", demo5)
	router.GET("/demo_6", demo6)
	router.GET("/demo_7", demo7)
	router.GET("/demo_8", demo8)
	router.GET("/demo_9", demo9)
	router.GET("/demo_10", demo10)
	router.POST("/password", checkPassword)
	router.POST("/restart", restartRpi)
	router.POST("/check_cable", checkCable)
	router.POST("/stop_stream", stopStream)
	router.POST("/shutdown", shutdownRpi)
	router.POST("/dhcp", changeToDhcp)
	router.POST("/server", changeServerAddress)
	router.POST("/static", changeToStatic)
	router.ServeFiles("/font/*filepath", http.Dir("font"))
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/pdf/*filepath", http.Dir("pdf"))
	router.Handler("GET", "/networkdata", networkDataStreamer)
	go StreamNetworkData(networkDataStreamer)
	_ = http.ListenAndServe(":9999", router)
}

func setupRemotePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled, _, _, mac := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	tmpl := template.Must(template.ParseFiles("html/setup-remote.html"))
	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
		Dhcp:            dhcpEnabled,
		Mac:             mac,
		DhcpChecked:     "",
		Version:         version,
	}
	if strings.Contains(dhcpEnabled, "yes") {
		data.DhcpChecked = "checked"
	}
	_ = tmpl.Execute(w, data)
}

func checkCable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _, _, _, active, _, _ := GetNetworkData()
	var responseData ChangeOutput
	responseData.Result = active
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseData)
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

func StreamNetworkData(streamer *sse.Streamer) {
	for {
		if streamCanRun {
			fmt.Println("streaming data")
			activeColor := "red"
			serverActiveColor := "red"
			interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled, active, _, mac := GetNetworkData()
			interfaceServerIpAddress := LoadSettingsFromConfigFile()
			serverAccessible := CheckServerIpAddress(interfaceServerIpAddress)
			if dhcpEnabled == "yes" {
				dhcpEnabled = "yes"
			} else {
				dhcpEnabled = "no"
			}
			serverActive := "server not accessible"
			if serverAccessible {
				serverActive = "server accessible"
				serverActiveColor = "green"
			}
			if active == "cable connected" {
				activeColor = "green"
			}
			streamer.SendString("", "networkdata", interfaceIpAddress+";"+interfaceMask+";"+interfaceGateway+";"+dhcpEnabled+";"+interfaceServerIpAddress+";"+mac+";"+active+";"+serverActive+";"+activeColor+";"+serverActiveColor+";"+mac)
		}
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
