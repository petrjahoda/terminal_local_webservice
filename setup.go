package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ChangeNetwork(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_ = r.ParseForm()
	ipaddress := r.Form["ipaddress"]
	gateway := r.Form["gateway"]
	mask := r.Form["mask"]
	serveripaddress := r.Form["serveripaddress"]
	pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if pattern.MatchString(ipaddress[0]) && pattern.MatchString(gateway[0]) {
		result, err := exec.Command("Powershell.exe", "netsh interface ipv4 set address name=\"Ethernet\" static "+ipaddress[0]+" "+mask[0]+" "+gateway[0]).Output()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		LogInfo("MAIN", "Change to static ip with result: "+string(result))
	}
	if pattern.MatchString(serveripaddress[0]) {
		configDirectory := filepath.Join(".", "config")
		configFileName := "config.json"
		configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
		data := ServerIpAddress{
			ServerIpAddress: serveripaddress[0],
		}
		file, _ := json.MarshalIndent(data, "", "  ")
		writingError := ioutil.WriteFile(configFullPath, file, 0666)
		LogInfo("MAIN", "Updating server ip address")
		if writingError != nil {
			LogError("MAIN", "Unable to update server ip address: "+writingError.Error())
		} else {
			LogInfo("MAIN", "Server ip address updated")
		}
	}
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	interfaceIpAddress, interfaceMask, interfaceGateway, interfaceDhcp := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
		Dhcp:            interfaceDhcp,
	}
	_ = tmpl.Execute(w, data)
}

func ChangeNetworkToDhcp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	println("Network is changing to DHCP")
	result, err := exec.Command("Powershell.exe", "netsh interface ipv4 set address name=\"Ethernet\" source=dhcp").Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("MAIN", "Changed to DHCP with result: "+string(result))

	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	interfaceIpAddress, interfaceMask, interfaceGateway, interfaceDhcp := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
		Dhcp:            interfaceDhcp,
	}
	_ = tmpl.Execute(w, data)
}
