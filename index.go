package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type ServerIpAddress struct {
	ServerIpAddress string
}

var HomepageLoaded bool

type HomepageData struct {
	IpAddress       string
	Mask            string
	Gateway         string
	ServerIpAddress string
	Dhcp            string
	Version         string
}

func GetNetworkData() (string, string, string, string) {
	LogInfo("STREAM", "Getting network data")
	var interfaceIpAddress string
	var notConnectedInterfaceIpAddress string
	var interfaceMask string
	var interfaceGateway string
	var notConnectedGateway string
	var interfaceDhcp string
	data, err := exec.Command("nmcli", "con", "show", "Wired connection 1").Output()
	if err != nil {
		LogError("STREAM", err.Error())
	}
	result := string(data)
	for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {
		if strings.Contains(line, "IP4.ADDRESS") {
			LogInfo("STREAM", "Processing: "+line)
			interfaceIpAddress = line[38:]
			LogInfo("STREAM", "Ip Address: "+interfaceIpAddress)
			interfaceIpAddress = interfaceIpAddress[:]
			splittedIpAddress := strings.Split(interfaceIpAddress, "/")
			maskNumber := splittedIpAddress[1]
			interfaceMask = CalculateMaskFrom(maskNumber)
			LogInfo("STREAM", "Mask: "+interfaceMask)
		} else if strings.Contains(line, "ipv4.addresses") {
			LogInfo("STREAM", "Processing: "+line)
			interfaceIpAddress = line[38:]
			LogInfo("STREAM", "Interface: "+interfaceIpAddress)
			if interfaceIpAddress != "  --" {
				notConnectedInterfaceIpAddress = interfaceIpAddress[:]
				LogInfo("STREAM", "Ip Address: "+notConnectedInterfaceIpAddress)
				splittedIpAddress := strings.Split(notConnectedInterfaceIpAddress, "/")
				maskNumber := splittedIpAddress[1]
				interfaceMask = CalculateMaskFrom(maskNumber)
				LogInfo("STREAM", "Mask: "+interfaceMask)
			}
		}
		if strings.Contains(line, "IP4.GATEWAY") {
			LogInfo("STREAM", "Processing: "+line)
			interfaceGateway = line[40:]
			LogInfo("STREAM", "Gateway: "+interfaceGateway)
			interfaceGateway = interfaceGateway[:]
		}
		if strings.Contains(line, "ipv4.method") {
			LogInfo("STREAM", "Processing: "+line)
			interfaceDhcp = line[40:]
			if strings.Contains(interfaceDhcp, "auto") {
				interfaceDhcp = "yes"
			} else {
				interfaceDhcp = "no"
			}
		}
		if strings.Contains(line, "ipv4.gateway") {
			LogInfo("STREAM", "Processing: "+line)
			notConnectedGateway = line[40:]
			LogInfo("STREAM", "Gateway: "+notConnectedGateway)
			notConnectedGateway = notConnectedGateway[:]
		}

	}
	if interfaceGateway == "" {
		interfaceGateway = notConnectedGateway + " not connected"
		interfaceIpAddress = notConnectedInterfaceIpAddress + " not connected"
		interfaceMask = interfaceMask + " not connected"
		interfaceDhcp = interfaceDhcp + " not connected"
	}
	return interfaceIpAddress, interfaceMask, interfaceGateway, interfaceDhcp
}

func Restart(http.ResponseWriter, *http.Request, httprouter.Params) {
	LogInfo("MAIN", "Restarting")
	start := time.Now()
	data, err := exec.Command("reboot").Output()
	if err != nil {
		LogError("MAIN", err.Error())
	}
	LogInfo("MAIN", "Restarted in "+time.Since(start).String()+" with result: "+string(data))
}

func Shutdown(http.ResponseWriter, *http.Request, httprouter.Params) {
	LogInfo("MAIN", "Shutting down")
	start := time.Now()
	data, err := exec.Command("poweroff").Output()
	if err != nil {
		LogError("MAIN", err.Error())
	}
	LogInfo("MAIN", "Shut down in "+time.Since(start).String()+" with result: "+string(data))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Index page Loading")
	start := time.Now()
	_ = r.ParseForm()
	interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	serverAccessible, url, interfaceServerIpAddress := CheckServerIpAddress(interfaceServerIpAddress)
	tmpl := template.Must(template.ParseFiles("html/index.html"))
	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		Dhcp:            dhcpEnabled,
		ServerIpAddress: url + " online",
		Version:         version,
	}
	if !serverAccessible {
		data.ServerIpAddress = url + " offline"
	}
	HomepageLoaded = true
	_ = tmpl.Execute(w, data)
	LogInfo("MAIN", "Homepage loaded in "+time.Since(start).String())
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
