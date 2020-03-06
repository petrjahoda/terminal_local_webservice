package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/sse"
	"github.com/kardianos/service"
	"html/template"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const version = "2020.1.3.6"
const programName = "Terminal local webservice"
const programDesription = "Display local web for asus terminals"
const deleteLogsAfter = 240 * time.Hour

type Page struct {
	Title string
	Body  []byte
}

type program struct{}

func (p *program) Start(s service.Service) error {
	LogInfo("MAIN", "Starting "+programName+" on "+s.Platform())
	go p.run()
	return nil
}

func (p *program) run() {
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
	_ = http.ListenAndServe(":80", router)
}

func (p *program) Stop(s service.Service) error {
	LogInfo("MAIN", "Stopped on platform "+s.Platform())
	return nil
}

func main() {
	serviceConfig := &service.Config{
		Name:        programName,
		DisplayName: programName,
		Description: programDesription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		LogError("MAIN", err.Error())
	}
	err = s.Run()
	if err != nil {
		LogError("MAIN", "Problem starting "+serviceConfig.Name)
	}
}

func StreamNetworkData(streamer *sse.Streamer) {
	timing := 20
	timeToSend := "20"
	refreshDone := true
	for {
		LogInfo("NETWORKDATA", "Streaming network data")
		start := time.Now()
		interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
		interfaceServerIpAddress := LoadSettingsFromConfigFile()
		serverAccessible, url, interfaceServerIpAddress := CheckServerIpAddress(interfaceServerIpAddress)
		if serverAccessible && !HomepageLoaded {
			timing = 0
			timeToSend = strconv.Itoa(timing)
			url = "http://localhost"
		} else if serverAccessible && HomepageLoaded {
			timing--
			timeToSend = strconv.Itoa(timing)
			refreshDone = false
		} else if !HomepageLoaded {
			timing = 20
			timeToSend = strconv.Itoa(timing)
		} else if !serverAccessible {
			if !refreshDone {
				timing = 0
				url = "http://localhost"
				timeToSend = strconv.Itoa(timing)
				refreshDone = true
			} else {
				timing = 20
				timeToSend = strconv.Itoa(timing)
			}
		}
		if timing < 0 {
			timing = 20
			timeToSend = strconv.Itoa(timing)
		}
		streamer.SendString("", "networkdata", interfaceIpAddress+";"+interfaceMask+";"+interfaceGateway+";"+dhcpEnabled+";"+timeToSend+";"+url+";"+interfaceServerIpAddress)
		LogInfo("NETWORKDATA", "Stream done in "+time.Since(start).String())
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
	LogInfo("NETWORKDATA", "Getting network data")
	var interfaceIpAddress string
	var interfaceMask string
	var interfaceGateway string
	var interfaceDhcp string
	if runtime.GOOS == "linux" {
		data, err := exec.Command("nmcli", "con", "show", "Wired connection 1").Output()
		if err != nil {
			LogError("NETWORKDATA", err.Error())
		}
		result := string(data)
		for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {
			if strings.Contains(line, "IP4.ADDRESS") {
				interfaceIpAddress = line[38:]
				LogInfo("NETWORKDATA", "Ip Address: "+interfaceIpAddress)
				interfaceIpAddress = interfaceIpAddress[:]
				splittedIpAddress := strings.Split(interfaceIpAddress, "/")
				maskNumber := splittedIpAddress[1]
				interfaceMask = CalculateMaskFrom(maskNumber)
				LogInfo("NETWORKDATA", "Mask: "+interfaceMask)
			}
			if strings.Contains(line, "IP4.GATEWAY") {
				interfaceGateway = line[40:]
				LogInfo("NETWORKDATA", "Gateway: "+interfaceGateway)
				interfaceGateway = interfaceGateway[:]
			}
			if strings.Contains(line, "ipv4.method") {
				interfaceDhcp = line[40:]
				if strings.Contains(interfaceDhcp, "auto") {
					interfaceDhcp = "yes"
				} else {
					interfaceDhcp = "no"
				}
			}
			if strings.Contains(line, "Wireless") {
				break
			}
		}
	} else {
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
	}
	if interfaceGateway == "" {
		interfaceGateway = "not connected"
		interfaceIpAddress = "not connected"
		interfaceMask = "not connected"
		interfaceDhcp = "not connected"
	}
	return interfaceIpAddress, interfaceMask, interfaceGateway, interfaceDhcp
}

func CalculateMaskFrom(maskNumber string) string {
	switch maskNumber {
	case "1":
		return "128.0.0.0"
	case "2":
		return "192.0.0.0"
	case "3":
		return "224.0.0.0"
	case "4":
		return "240.0.0.0"
	case "5":
		return "248.0.0.0"
	case "6":
		return "252.0.0.0"
	case "7":
		return "254.0.0.0"
	case "8":
		return "255.0.0.0"
	case "9":
		return "255.128.0.0"
	case "10":
		return "255.192.0.0"
	case "11":
		return "255.224.0.0"
	case "12":
		return "255.240.0.0"
	case "13":
		return "255.248.0.0"
	case "14":
		return "255.252.0.0"
	case "15":
		return "255.254.0.0"
	case "16":
		return "255.255.0.0"
	case "17":
		return "255.255.128.0"
	case "18":
		return "255.255.192.0"
	case "19":
		return "255.255.224.0"
	case "20":
		return "255.255.240.0"
	case "21":
		return "255.255.248.0"
	case "22":
		return "255.255.252.0"
	case "23":
		return "255.255.254.0"
	case "24":
		return "255.255.255.0"
	case "25":
		return "255.255.255.128"
	case "26":
		return "255.255.255.192"
	case "27":
		return "255.255.255.224"
	case "28":
		return "255.255.255.240"
	case "29":
		return "255.255.255.248"
	case "30":
		return "255.255.255.252"
	case "31":
		return "255.255.255.254"
	case "32":
		return "255.255.255.255"
	}
	return "-"
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("html/" + tmpl + ".html")
	_ = t.Execute(w, p)
}

func darcula(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "css/darcula.css")
}

func metrojs(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "js/metro.min.js")
}

func metrocss(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "css/metro-all.css")
}

func image(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	http.ServeFile(writer, request, "image.png")
}
