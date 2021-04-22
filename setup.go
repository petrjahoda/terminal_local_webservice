package main

import (
	"encoding/json"
	"fmt"
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

type ChangeInput struct {
	Password  string
	IpAddress string
	Mask      string
	Gateway   string
	Server    string
}

type ChangeOutput struct {
	Result string
}

func setupPage(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	streamSync.Lock()
	streamCanRun = false
	streamSync.Unlock()
	interfaceIpAddress, interfaceMask, interfaceGateway, dhcpEnabled := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	tmpl := template.Must(template.ParseFiles("html/setup.html"))
	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
		Dhcp:            dhcpEnabled,
		DhcpChecked:     "",
		Version:         version,
	}
	if strings.Contains(dhcpEnabled, "yes") {
		data.DhcpChecked = "checked"
	}
	_ = tmpl.Execute(w, data)
}

func changeToStatic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data ChangeInput
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		var responseData ChangeOutput
		responseData.Result = "nok"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseData)
		return
	}
	if data.Password == "3600" {
		pattern := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
		if pattern.MatchString(data.IpAddress) && pattern.MatchString(data.Gateway) {
			maskNumber := GetMaskNumberFrom(data.Mask)
			result, err := exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "manual", "ipv4.addresses", data.IpAddress+"/"+maskNumber, "ipv4.gateway", data.Gateway).Output()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(result)
			result, err = exec.Command("nmcli", "con", "up", "Wired connection 1").Output()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(result)
		}
		if len(data.Server) > 0 {
			configDirectory := filepath.Join(".", "config")
			configFileName := "config.json"
			configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
			data := ServerIpAddress{
				ServerIpAddress: data.Server,
			}
			file, _ := json.MarshalIndent(data, "", "  ")
			_ = ioutil.WriteFile(configFullPath, file, 0666)
		}
		var responseData ChangeOutput
		responseData.Result = "ok"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseData)
		return
	}
	var responseData PasswordOutput
	responseData.Result = "nok"
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseData)
}

func changeToDhcp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data ChangeInput
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		var responseData ChangeOutput
		responseData.Result = "nok"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseData)
		return
	}
	fmt.Println(data)
	if data.Password == "3600" {
		result, err := exec.Command("nmcli", "con", "mod", "Wired connection 1", "ipv4.method", "auto").Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(result)
		if len(data.Server) > 0 {
			configDirectory := filepath.Join(".", "config")
			configFileName := "config.json"
			configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
			data := ServerIpAddress{
				ServerIpAddress: data.Server,
			}
			file, _ := json.MarshalIndent(data, "", "  ")
			_ = ioutil.WriteFile(configFullPath, file, 0666)
		}
		var responseData ChangeOutput
		responseData.Result = "ok"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseData)
		return
	}
	var responseData ChangeOutput
	responseData.Result = "nok"
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseData)
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

func CheckServerIpAddress(interfaceServerIpAddress string) bool {
	seconds := 2
	timeOut := time.Duration(seconds) * time.Second
	fmt.Println(interfaceServerIpAddress)
	result, err := net.DialTimeout("tcp", interfaceServerIpAddress, timeOut)
	fmt.Println(result)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
