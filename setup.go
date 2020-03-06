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
	"runtime"
	"strings"
	"time"
)

func ChangeNetwork(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Change network Loading")
	start := time.Now()
	_ = r.ParseForm()
	ipaddress := r.Form["ipaddress"]
	gateway := r.Form["gateway"]
	mask := r.Form["mask"]
	serveripaddress := r.Form["serveripaddress"]
	pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	LogInfo("MAIN", "New IP Address: "+ipaddress[0])
	LogInfo("MAIN", "New Mask: "+mask[0])
	LogInfo("MAIN", "New Gateway: "+gateway[0])
	if pattern.MatchString(ipaddress[0]) && pattern.MatchString(gateway[0]) {
		if runtime.GOOS == "linux" {
			maskNumber := GetMaskNumberFrom(mask[0])
			LogInfo("MAIN", "New Mask Number: "+maskNumber)
			result, err := exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "manual", "ipv4.addresses", ipaddress[0]+"/"+maskNumber, "ipv4.gateway", gateway[0]).Output()
			if err != nil {
				LogError("MAIN", err.Error())
			}

			LogInfo("MAIN", string(result))
			result, err = exec.Command("nmcli", "con", "up", "Wired connection 1").Output()
			if err != nil {
				LogError("MAIN", err.Error())
			}
			LogInfo("MAIN", "Change to static ip with result: "+string(result))
		} else {
			result, err := exec.Command("Powershell.exe", "netsh interface ipv4 set address name=\"Ethernet\" static "+ipaddress[0]+" "+mask[0]+" "+gateway[0]).Output()
			if err != nil {
				LogError("MAIN", err.Error())
			}
			LogInfo("MAIN", "Change to static ip with result: "+string(result))
		}
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
		LogInfo("MAIN", "Updating server ip address")
		if writingError != nil {
			LogError("MAIN", "Unable to update server ip address: "+writingError.Error())
		} else {
			LogInfo("MAIN", "Server ip address updated")
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
	LogInfo("MAIN", "Change network loaded in "+time.Since(start).String())

}

func GetMaskNumberFrom(maskNumber string) string {
	switch maskNumber {
	case "128.0.0.0":
		return "1"
	case "192.0.0.0":
		return "2"
	case "224.0.0.0":
		return "3"
	case "240.0.0.0":
		return "4"
	case "248.0.0.0":
		return "5"
	case "252.0.0.0":
		return "6"
	case "254.0.0.0":
		return "7"
	case "255.0.0.0":
		return "8"
	case "255.128.0.0":
		return "9"
	case "255.192.0.0":
		return "10"
	case "255.224.0.0":
		return "11"
	case "255.240.0.0":
		return "12"
	case "255.248.0.0":
		return "13"
	case "255.252.0.0":
		return "14"
	case "255.254.0.0":
		return "15"
	case "255.255.0.0":
		return "16"
	case "255.255.128.0":
		return "17"
	case "255.255.192.0":
		return "18"
	case "255.255.224.0":
		return "19"
	case "255.255.240.0":
		return "20"
	case "255.255.248.0":
		return "21"
	case "255.255.252.0":
		return "22"
	case "255.255.254.0":
		return "23"
	case "255.255.255.0":
		return "24"
	case "255.255.255.128":
		return "25"
	case "255.255.255.192":
		return "26"
	case "255.255.255.224":
		return "27"
	case "255.255.255.240":
		return "28"
	case "255.255.255.248":
		return "29"
	case "255.255.255.252":
		return "30"
	case "255.255.255.254":
		return "31"
	case "255.255.255.255":
		return "32"
	}
	return "0"
}

func ChangeNetworkToDhcp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Change DHCP loading")
	start := time.Now()
	if runtime.GOOS == "linux" {
		result, err := exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "auto").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Changed to DHCP with result: "+string(result))
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
		LogInfo("MAIN", "Loaded in "+time.Since(start).String())
	} else {
		result, err := exec.Command("Powershell.exe", "netsh interface ipv4 set address name=\"Ethernet\" source=dhcp").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Changed to DHCP with result: "+string(result))
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
		LogInfo("MAIN", "Change DHCP loaded in "+time.Since(start).String())
	}
}

func CheckServerIpAddress(interfaceServerIpAddress string) (bool, string, string) {
	LogInfo("STREAM", "Checking server ip address")
	start := time.Now()
	serverAccessible := false
	url := ""
	hostName := interfaceServerIpAddress
	portNum := "80"
	seconds := 2
	timeOut := time.Duration(seconds) * time.Second
	_, err := net.DialTimeout("tcp", hostName+":"+portNum, timeOut)
	if err != nil {
		LogError("STREAM", interfaceServerIpAddress+" not accessible: "+err.Error())
		interfaceServerIpAddress += " not accessible"
	} else {
		LogInfo("STREAM", interfaceServerIpAddress+" accessible")
		serverAccessible = true
		url = "http://" + interfaceServerIpAddress + "/"
	}
	LogInfo("STREAM", "Checked in "+time.Since(start).String())
	return serverAccessible, url, interfaceServerIpAddress
}
