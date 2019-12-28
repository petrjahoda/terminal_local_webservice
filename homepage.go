package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kbinani/screenshot"
	"html/template"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func Screenshot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("SCREENSHOT", "Loading")
	start := time.Now()
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		img, err := screenshot.CaptureDisplay(i)
		if err != nil {
			LogError("SCREENSHOT", "Error generating screenshot: "+err.Error())
			continue
		}
		fileName := "image.png"
		file, _ := os.Create(fileName)
		_ = png.Encode(file, img)
		LogInfo("SCREENSHOT", "Generated screenshot: "+fileName)
		file.Close()
	}
	HomepageLoaded = false
	renderTemplate(w, "screenshot", &Page{})
	LogInfo("SCREENSHOT", "Loaded in "+time.Since(start).String())
}

func RestartBrowser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("RESTART BROWSER", "Loading")
	start := time.Now()
	data, err := exec.Command("Powershell.exe", "Stop-Process -Name chrome").Output()
	if err != nil {
		LogError("RESTART BROWSER", err.Error())
	}
	data, err = exec.Command("Powershell.exe", "Start-Process \"chrome.exe\"", "\"--kiosk --disable-pinch --app=http://localhost:8000\"").Output()
	LogInfo("RESTART BROWSER", string(data))
	if err != nil {
		LogError("RESTART BROWSER", err.Error())
	}
	LogInfo("RESTART BROWSER", "Loaded in "+time.Since(start).String())
}

func Restart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("RESTART", "Loading")
	start := time.Now()
	data, err := exec.Command("Powershell.exe", "Restart-Computer").Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("RESTART", "Loaded in "+time.Since(start).String()+" with result: "+string(data))
}

func Shutdown(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("SHUTDOWN", "Loading")
	start := time.Now()
	data, err := exec.Command("Powershell.exe", "Stop-Computer").Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	LogInfo("SHUTDOWN", "Loaded in "+time.Since(start).String()+" with result: "+string(data))
}

func Homepage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("HOMEPAGE", "Loading")
	start := time.Now()
	_ = r.ParseForm()
	tmpl := template.Must(template.ParseFiles("html/homepage.html"))
	println(version)
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
	LogInfo("HOMEPAGE", "Loaded in "+time.Since(start).String())
}

func Setup(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("SETUP", "Loading")
	start := time.Now()
	_ = r.ParseForm()
	password := r.Form["password"]
	if password[0] == "2011" {
		HomepageLoaded = false
		renderTemplate(w, "setup", &Page{})
	} else {
		LogInfo("SETUP", "Bad password")
		HomepageLoaded = true
		_ = r.ParseForm()
		tmpl := template.Must(template.ParseFiles("html/homepage.html"))
		LogInfo("SETUP", version)
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
		LogInfo("SETUP", "Loaded in "+time.Since(start).String())
	}
}
func Password(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	LogInfo("PASSWORD", "Loading")
	start := time.Now()
	HomepageLoaded = false
	renderTemplate(w, "password", &Page{})
	LogInfo("PASSWORD", "Loaded in "+time.Since(start).String())
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
