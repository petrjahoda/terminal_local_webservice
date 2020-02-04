package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"html/template"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var Interface = "Ethernet"

const version = "2020.1.2.4"
const deleteLogsAfter = 240 * time.Hour

type Page struct {
	Title string
	Body  []byte
}

func main() {
	LogDirectoryFileCheck("MAIN")
	CreateConfigIfNotExists()

	router := httprouter.New()
	timeStreamer := sse.New()
	networkDataStreamer := sse.New()

	router.GET("/", Homepage)
	router.GET("/screenshot", Screenshot)
	router.GET("/password", Password)
	router.GET("/changenetwork", ChangeNetwork)
	router.GET("/changenetworktodhcp", ChangeNetworkToDhcp)
	router.GET("/restart", Restart)
	router.GET("/restartbrowser", RestartBrowser)
	router.GET("/shutdown", Shutdown)
	router.GET("/setup", Setup)
	router.GET("/css/darcula.css", darcula)
	router.GET("/js/metro.min.js", metrojs)
	router.GET("/css/metro-all.css", metrocss)
	router.GET("/image.png", image)

	router.Handler("GET", "/listen", timeStreamer)
	router.Handler("GET", "/networkdata", networkDataStreamer)
	go StreamTime(timeStreamer)
	go StreamNetworkData(networkDataStreamer)
	LogInfo("MAIN", "Server running")
	_ = http.ListenAndServe(":8000", router)
}

func StreamNetworkData(streamer *sse.Streamer) {
	timing := 20
	timeToSend := "20"
	for {
		interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
		interfaceServerIpAddress := LoadSettingsFromConfigFile()
		serverAccessible, url, interfaceServerIpAddress := CheckServerIpAddress(interfaceServerIpAddress)
		if serverAccessible && HomepageLoaded {
			timing--
			timeToSend = strconv.Itoa(timing)
		}
		if !serverAccessible {
			timing = 20
			timeToSend = strconv.Itoa(timing)
		}
		if timing < 0 {
			timing = 20
			timeToSend = strconv.Itoa(timing)
		}
		if !HomepageLoaded {
			timing = 20
		}
		streamer.SendString("", "networkdata", interfaceIpAddress+";"+interfaceMask+";"+interfaceGateway+";"+dhcpEnabled+";"+timeToSend+";"+url+";"+interfaceServerIpAddress)
		time.Sleep(1 * time.Second)
	}
}

func StreamTime(streamer *sse.Streamer) {
	for {
		streamer.SendString("", "time", time.Now().Format("15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

func GetNetworkData() (string, string, string, string) {
	var interfaceIpAddress string
	var interfaceMask string
	var interfaceGateway string
	var interfaceDhcp string

	data, err := exec.Command("Powershell.exe", "ipconfig /all").Output()

	if err != nil {
		LogError("MAIN", err.Error())
	}
	result := string(data)
	ethernetStarts := false
	for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {
		if strings.Contains(line, "Ethernet") {
			ethernetStarts = true
		}
		if ethernetStarts {
			if strings.Contains(line, "IPv4 Address") {
				interfaceIpAddress = line[38:]
				interfaceIpAddress = interfaceIpAddress[:len(interfaceIpAddress)-1]

			}
			if strings.Contains(line, "Subnet Mask") {
				interfaceMask = line[38:]
				interfaceMask = interfaceMask[:len(interfaceMask)-1]

			}
			if strings.Contains(line, "Default Gateway") {
				interfaceGateway = line[38:]
				interfaceGateway = interfaceGateway[:len(interfaceGateway)-1]

			}
			if strings.Contains(line, "DHCP Enabled") {
				interfaceDhcp = line[38:]
				interfaceDhcp = interfaceDhcp[:len(interfaceDhcp)-1]
			}
			if strings.Contains(line, "Wireless") {
				break
			}

		}
	}
	if interfaceGateway == "" {
		interfaceGateway = "not connected"
		interfaceIpAddress = "not connected"
		interfaceMask = "not connected"
		interfaceDhcp = "not connected"
	}
	return interfaceIpAddress, interfaceMask, interfaceGateway, interfaceDhcp
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
