package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/route"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ServerIpAddress struct {
	IpAddress string
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
	if interfaceGateway == "" {
		interfaceGateway = "not connected"
	}
	return interfaceIpAddress, interfaceMask, interfaceGateway
}

func GetGateway() string {
	rib, _ := route.FetchRIB(0, route.RIBTypeRoute, 0)
	messages, err := route.ParseRIB(route.RIBTypeRoute, rib)
	if err != nil {
		println(err.Error())
	}
	for _, message := range messages {
		route_message := message.(*route.RouteMessage)
		addresses := route_message.Addrs

		var destination, gateway *route.Inet4Addr
		ok := false
		if destination, ok = addresses[0].(*route.Inet4Addr); !ok {
			continue
		}

		if gateway, ok = addresses[1].(*route.Inet4Addr); !ok {
			continue
		}

		if destination == nil || gateway == nil {
			continue
		}
		var defaultRoute = [4]byte{0, 0, 0, 0}
		if destination.IP == defaultRoute {
			return strconv.Itoa(int(gateway.IP[0])) + "." + strconv.Itoa(int(gateway.IP[1])) + "." + strconv.Itoa(int(gateway.IP[2])) + "." + strconv.Itoa(int(gateway.IP[3]))
		}
	}
	return "gateway not set"
}

func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		panic("ipv4Mask: len must be 4 bytes")
	}

	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

type HomepageData struct {
	IpAddress       string
	Mask            string
	Gateway         string
	ServerIpAddress string
}
