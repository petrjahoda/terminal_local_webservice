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
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const version = "2021.2.1.23"
const programName = "Terminal local webservice"
const programDescription = "Display local web for rpi terminals"

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
	initiateRaspberry()
	router := httprouter.New()
	networkDataStreamer := sse.New()
	router.GET("/image.png", image)
	router.GET("/", indexPage)
	router.GET("/screenshot", screenshotPage)
	router.GET("/setup", setupPage)
	router.POST("/password", checkPassword)
	router.POST("/restart", restartRpi)
	router.POST("/stop_stream", stopStream)
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

func initiateRaspberry() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	readFile, _ := ioutil.ReadFile(configFullPath)
	ConfigFile := ServerIpAddress{}
	_ = json.Unmarshal(readFile, &ConfigFile)
	ipaddress := ConfigFile.IpAddress
	mask := ConfigFile.Mask
	gateway := ConfigFile.Gateway
	dhcp := ConfigFile.Dhcp
	if dhcp == "true" {
		fmt.Println("INITIATE DHCP TRUE")
		result, err := exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "auto").Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("INITIATE DHCP TRUE RESULT: " + string(result))
	} else {
		fmt.Println("INITIATE DHCP FALSE")
		pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
		if pattern.MatchString(ipaddress) && pattern.MatchString(gateway) {
			maskNumber := GetMaskNumberFrom(mask)
			result, err := exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "manual", "ipv4.addresses", ipaddress+"/"+maskNumber, "ipv4.gateway", gateway).Output()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("INITIATE DHCP FALSE RESULT1 " + string(result))
			result, err = exec.Command("nmcli", "con", "up", "Wired connection 1").Output()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("INITIATE DHCP FALSE RESULT2 " + string(result))
		}
	}
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
			activColor := "red"
			serverActivColor := "red"
			interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled, active, result := GetNetworkData()
			interfaceServerIpAddress := LoadSettingsFromConfigFile()
			serverAccessible := CheckServerIpAddress(interfaceServerIpAddress)
			if dhcpEnabled == "yes" {
				dhcpEnabled = "ano"
			} else {
				dhcpEnabled = "ne"
			}
			serverActive := "server nedostupný"
			if serverAccessible {
				serverActive = "server dostupný"
				serverActivColor = "green"
			}
			if active == "kabel zapojený" {
				activColor = "green"
			}
			streamer.SendString("", "networkdata", interfaceIpAddress+";"+interfaceMask+";"+interfaceGateway+";"+dhcpEnabled+";"+interfaceServerIpAddress+";"+result+";"+active+";"+serverActive+";"+activColor+";"+serverActivColor)
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
