package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/kbinani/screenshot"
	"html/template"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func Screenshot(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Screenshot loading")
	start := time.Now()
	if runtime.GOOS == "linux" {
		data, err := exec.Command("/usr/bin/maim", "image.png").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Screenshot taken: "+string(data))
	} else {
		n := screenshot.NumActiveDisplays()
		for i := 0; i < n; i++ {
			img, err := screenshot.CaptureDisplay(i)
			if err != nil {
				LogError("MAIN", "Error generating screenshot: "+err.Error())
				continue
			}
			fileName := "image.png"
			file, _ := os.Create(fileName)
			_ = png.Encode(file, img)
			LogInfo("MAIN", "Generated screenshot: "+fileName)
			file.Close()
		}

	}
	HomepageLoaded = false
	renderTemplate(w, "screenshot", &Page{})
	LogInfo("MAIN", "Screenshot loaded in "+time.Since(start).String())
}

func Restart(http.ResponseWriter, *http.Request, httprouter.Params) {
	LogInfo("MAIN", "Restarting")
	start := time.Now()
	if runtime.GOOS == "linux" {
		data, err := exec.Command("reboot").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Restarted in "+time.Since(start).String()+" with result: "+string(data))
	} else {
		data, err := exec.Command("Powershell.exe", "Restart-Computer").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Restarted in "+time.Since(start).String()+" with result: "+string(data))
	}
}

func Shutdown(http.ResponseWriter, *http.Request, httprouter.Params) {
	LogInfo("MAIN", "Shutting down")
	start := time.Now()
	if runtime.GOOS == "linux" {
		data, err := exec.Command("poweroff").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Shut down in "+time.Since(start).String()+" with result: "+string(data))
	} else {
		data, err := exec.Command("Powershell.exe", "Stop-Computer").Output()
		if err != nil {
			LogError("MAIN", err.Error())
		}
		LogInfo("MAIN", "Shut down in "+time.Since(start).String()+" with result: "+string(data))
	}
}

func Homepage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Homepage Loading")
	start := time.Now()
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
	LogInfo("MAIN", "Homepage loaded in "+time.Since(start).String())
}

func Setup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Setup loading")
	start := time.Now()
	_ = r.ParseForm()
	password := r.Form["password"]
	println(len(password))
	if password[0] == "2011" {
		HomepageLoaded = false
		renderTemplate(w, "setup", &Page{})
	} else {
		LogInfo("MAIN", "Bad password")
		HomepageLoaded = true
		_ = r.ParseForm()
		tmpl := template.Must(template.ParseFiles("html/homepage.html"))
		LogInfo("MAIN", version)
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
		LogInfo("MAIN", "Setup loaded in "+time.Since(start).String())
	}
}
func Password(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	LogInfo("MAIN", "Password loading")
	start := time.Now()
	HomepageLoaded = false
	renderTemplate(w, "password", &Page{})
	LogInfo("MAIN", "Password loaded in "+time.Since(start).String())
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
				ServerIpAddress: "",
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
