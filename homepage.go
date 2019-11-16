package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ServerIpAddress struct {
	ServerIpAddress string
}

type HomepageData struct {
	IpAddress       string
	Mask            string
	Gateway         string
	ServerIpAddress string
	Timer           string
	Url             string
}

func Screenshot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Generating screenshot")
	data, err := exec.Command("Powershell.exe", "screenshot.exe").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("MAIN", string(data))
	//n := screenshot.NumActiveDisplays()
	//LogInfo("MAIN", "Displays: "+strconv.Itoa(n))
	//
	//for i := 0; i < n; i++ {
	//	bounds := screenshot.GetDisplayBounds(i)
	//	LogInfo("MAIN", "Bounds: "+bounds.String())
	//	img, err := screenshot.CaptureRect(bounds)
	//	if err != nil {
	//		LogError("MAIN", "Error generating screenshot: "+err.Error())
	//		continue
	//	}
	//	fileName := "image.png"
	//	file, _ := os.Create(fileName)
	//	defer file.Close()
	//	_ = png.Encode(file, img)
	//	LogInfo("MAIN", "Generated screenshot: "+fileName)
	//}
	LogInfo("MAIN", "Generating finished")
	renderTemplate(w, "screenshot", &Page{})
}

func Restart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data, err := exec.Command("Powershell.exe", "Restart-Computer").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("MAIN", string(data))
}

func Shutdown(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data, err := exec.Command("Powershell.exe", "Stop-Computer").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("MAIN", string(data))
}

func Homepage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))

	interfaceIpAddress, interfaceMask, interfaceGateway := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()
	timer := "86400"
	url := ""
	if interfaceIpAddress != "not connected" {
		println("timing")
		timer = "20"
		url = "http://" + interfaceServerIpAddress + "/"
	}

	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
		Timer:           timer,
		Url:             url,
	}
	_ = tmpl.Execute(w, data)
}

func Setup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/setup.html"))

	interfaceIpAddress, interfaceMask, interfaceGateway := GetNetworkData()
	interfaceServerIpAddress := LoadSettingsFromConfigFile()

	data := HomepageData{
		IpAddress:       interfaceIpAddress,
		Mask:            interfaceMask,
		Gateway:         interfaceGateway,
		ServerIpAddress: interfaceServerIpAddress,
	}
	_ = tmpl.Execute(w, data)
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

func CreateConfigIfNotExists() {
	configDirectory := filepath.Join(".", "config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")

	if _, checkPathError := os.Stat(configFullPath); checkPathError == nil {
		LogInfo("MAIN", "Config file already exists")
	} else if os.IsNotExist(checkPathError) {
		LogWarning("MAIN", "Config file does not exist, creating")
		mkdirError := os.MkdirAll(configDirectory, 0777)
		if mkdirError != nil {
			LogError("MAIN", "Unable to create directory for config file: "+mkdirError.Error())
		} else {
			LogInfo("MAIN", "Directory for config file created")
			data := ServerIpAddress{
				ServerIpAddress: "192.168.1.11",
			}
			file, _ := json.MarshalIndent(data, "", "  ")
			writingError := ioutil.WriteFile(configFullPath, file, 0666)
			LogInfo("MAIN", "Writing data to JSON file")
			if writingError != nil {
				LogError("MAIN", "Unable to write data to config file: "+writingError.Error())
			} else {
				LogInfo("MAIN", "Data written to config file")
			}
		}
	} else {
		LogError("MAIN", "Config file does not exist")
	}
}

func GetNetworkData() (string, string, string) {
	var interfaceIpAddress string
	var interfaceMask string
	var interfaceGateway string

	data, err := exec.Command("Powershell.exe", "ipconfig").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	result := string(data)
	println(result)
	for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {
		if strings.Contains(line, "IPv4 Address") {
			interfaceIpAddress = line[38:]
		}
		if strings.Contains(line, "Subnet Mask") {
			interfaceMask = line[38:]
		}
		if strings.Contains(line, "Default Gateway") {
			interfaceGateway = line[38:]
		}
	}
	if interfaceGateway == "" {
		interfaceGateway = "not connected"

	}
	if interfaceIpAddress == "" {
		interfaceIpAddress = "not connected"

	}
	if interfaceMask == "" {
		interfaceMask = "not connected"
	}

	return interfaceIpAddress, interfaceMask, interfaceGateway
}
