package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kbinani/screenshot"
	"html/template"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type ServerIpAddress struct {
	IpAddress string
}

func Screenshot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("MAIN", "Generating screenshot")
	n := screenshot.NumActiveDisplays()
	LogInfo("MAIN", "Displays: "+strconv.Itoa(n))

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		LogInfo("MAIN", "Bounds: "+bounds.String())
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			LogError("MAIN", "Error generating screenshot: "+err.Error())
			continue
		}
		fileName := "image.png"
		file, _ := os.Create(fileName)
		defer file.Close()
		_ = png.Encode(file, img)
		LogInfo("MAIN", "Generated screenshot: "+fileName)
	}
	LogInfo("MAIN", "Generating finished")
	renderTemplate(w, "screenshot", &Page{})

}

func Homepage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("homepage.html"))

	interfaces, _ := net.Interfaces()
	interfaceIpAddress, interfaceMask, interfaceGateway := GetNetworkData(interfaces)

	CreateConfigIfNotExists()
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

	configDirectory := filepath.Join(".", "Config")
	configFileName := "config.json"
	configFullPath := strings.Join([]string{configDirectory, configFileName}, "/")
	readFile, _ := ioutil.ReadFile(configFullPath)
	ConfigFile := ServerIpAddress{}
	_ = json.Unmarshal(readFile, &ConfigFile)
	ServerIpAddress := ConfigFile.IpAddress
	return ServerIpAddress
}

func CreateConfigIfNotExists() {
	configDirectory := filepath.Join(".", "Config")
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
				IpAddress: "192.168.1.11",
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

func GetNetworkData(interfaces []net.Interface) (string, string, string) {
	var interfaceIpAddress string
	var interfaceMask string
	var interfaceGateway string
	for _, requestedInterface := range interfaces {
		if requestedInterface.Name == Interface {
			addrs, err := requestedInterface.Addrs()
			if err != nil {
				println("Bad interface")
			}
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP.To4()
					if ip != nil {
						interfaceIpAddress = ip.String()
						mask := ip.DefaultMask()
						interfaceMask = fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
						interfaceGateway = GetGateway()

					}
				case *net.IPAddr:
					ip = v.IP
					if ip != nil {
						interfaceIpAddress = ip.String()
						mask := ip.DefaultMask()
						interfaceMask = fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
						interfaceGateway = GetGateway()
					}
				}
			}
		}
	}

	if interfaceIpAddress == "" {
		interfaceIpAddress = "not connected"

	}
	if interfaceMask == "" {
		interfaceMask = "not connected"
	}

	return interfaceIpAddress, interfaceMask, interfaceGateway
}

type HomepageData struct {
	IpAddress       string
	Mask            string
	Gateway         string
	ServerIpAddress string
}

func GetGateway() string {
	data, err := exec.Command("Powershell.exe", "Get-NetIPConfiguration -InterfaceIndex 15").Output()

	if err != nil {
		fmt.Println("Error: ", err)
	}
	result := string(data)
	for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {
		if strings.Contains(line, "IPv4DefaultGateway") {
			return line[22:]
		}
	}
	return "not connected"
}
