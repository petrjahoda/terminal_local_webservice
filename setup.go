package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func ChangeNetwork(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("NETWORK CHANGE", "Loading")
	start := time.Now()
	_ = r.ParseForm()
	ipaddress := r.Form["ipaddress"]
	gateway := r.Form["gateway"]
	mask := r.Form["mask"]
	serveripaddress := r.Form["serveripaddress"]
	pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if pattern.MatchString(ipaddress[0]) && pattern.MatchString(gateway[0]) {
		result, err := exec.Command("Powershell.exe", "netsh interface ipv4 set address name=\"Ethernet\" static "+ipaddress[0]+" "+mask[0]+" "+gateway[0]).Output()
		if err != nil {
			LogError("NETWORK CHANGE", err.Error())
		}
		LogInfo("NETWORK CHANGE", "Change to static ip with result: "+string(result))
	}
	if len(serveripaddress[0]) > 0 {
		configDirectory := filepath.Join(".", "config")
		configFileName := "config.json"
		configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
		data := ServerIpAddress{
			ServerIpAddress: serveripaddress[0],
		}
		file, _ := json.MarshalIndent(data, "", "  ")
		writingError := ioutil.WriteFile(configFullPath, file, 0666)
		LogInfo("NETWORK CHANGE", "Updating server ip address")
		if writingError != nil {
			LogError("NETWORK CHANGE", "Unable to update server ip address: "+writingError.Error())
		} else {
			LogInfo("NETWORK CHANGE", "Server ip address updated")
		}
	}
	HomepageLoaded = true
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	data := HomepageData{
		IpAddress:       "",
		Mask:            "",
		Gateway:         "",
		ServerIpAddress: "",
		Dhcp:            "",
		Version:         version,
	}
	HomepageLoaded = true
	_ = tmpl.Execute(w, data)
	LogInfo("NETWORK CHANGE", "Loaded in "+time.Since(start).String())

}

func ChangeNetworkToDhcp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("NETWORK CHANGE TO DHCP", "Loading")
	start := time.Now()
	result, err := exec.Command("Powershell.exe", "netsh interface ipv4 set address name=\"Ethernet\" source=dhcp").Output()
	if err != nil {
		LogError("NETWORK CHANGE TO DHCP", err.Error())
	}
	LogInfo("NETWORK CHANGE TO DHCP", "Changed to DHCP with result: "+string(result))
	HomepageLoaded = true
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	data := HomepageData{
		IpAddress:       "",
		Mask:            "",
		Gateway:         "",
		ServerIpAddress: "",
		Dhcp:            "",
		Version:         version,
	}
	HomepageLoaded = true
	_ = tmpl.Execute(w, data)
	LogInfo("NETWORK CHANGE TO DHCP", "Loaded in "+time.Since(start).String())
}

func CheckServerIpAddress(interfaceServerIpAddress string) (bool, string, string) {
	serverAccessible := false
	url := ""
	hostName := interfaceServerIpAddress
	portNum := "80"
	seconds := 2
	timeOut := time.Duration(seconds) * time.Second
	_, err := net.DialTimeout("tcp", hostName+":"+portNum, timeOut)
	if err != nil {
		LogError("MAIN", interfaceServerIpAddress+" not accessible: "+err.Error())
		interfaceServerIpAddress += " not accessible"
	} else {
		LogInfo("MAIN", interfaceServerIpAddress+" accessible")
		serverAccessible = true
		url = "http://" + interfaceServerIpAddress + "/"
	}
	return serverAccessible, url, interfaceServerIpAddress
}
