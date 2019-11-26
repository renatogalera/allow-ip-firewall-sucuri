package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	externalip "github.com/GlenDC/go-external-ip"
	"github.com/joho/godotenv"
)

type SucuriSettings struct {
	Output      struct {
		Domain                  string        `json:"domain"`
		WhitelistList           []string      `json:"whitelist_list"`
		BlacklistList           []interface{} `json:"blacklist_list"`
		InternalDomainIP        string        `json:"internal_domain_ip"`
		InternalDomainDebugList []string      `json:"internal_domain_debug_list"`
	} `json:"output"`
	Verbose int `json:"verbose"`
}

var API_KEY string
var API_SECRET string
var NEWIPADDR string

func argParse() error {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	API_KEY = os.Getenv("API_KEY")
	if API_KEY == "" {
		log.Fatal("Need to define API_KEY var")
	}
	API_SECRET = os.Getenv("API_SECRET")
	if API_SECRET == "" {
		log.Fatal("Need to define API_SECRET var")
	}
	return nil
}

func getMyIP() string {
	consensus := externalip.DefaultConsensus(nil, nil)
	currentIP, err := consensus.ExternalIP()
	if err != nil {
		log.Println("Error collecting external IP", err)
	}
	TARGETIP := currentIP.String()
	return TARGETIP
}

func getAllowIP() SucuriSettings {
	var s SucuriSettings
	client := &http.Client{
		Timeout:time.Second * 10,
	}
	req, err := client.Get("https://waf.sucuri.net/api?v2&k=" + API_KEY + "&s=" + API_SECRET + "&a=show_settings")
	if err != nil {
		log.Fatal("Error to get info", err.Error())
	}
	defer req.Body.Close()
	if req.StatusCode == 200 {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal("Error to read sucuri page", err.Error())
		}
		json.Unmarshal(body, &s)
		return s
	}
	return s
}

func addIP()  {
	client := &http.Client{
		Timeout:time.Second * 10,
	}
	req, err := client.Get("https://waf.sucuri.net/api?k=" + API_KEY + "&s=" + API_SECRET + "&a=whitelist")
	if err != nil {
		log.Fatal("Error to get info", err.Error())
	}
	defer req.Body.Close()
	if req.StatusCode == 200 {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal("Error to read sucuri page", err.Error())
		}
		fmt.Println(string(body))
	}
}

func checkAllowIP(myactip string) bool {
	s := getAllowIP()
	for _, v := range s.Output.WhitelistList {
		if v == myactip {
			return true
		}
	}
	return false
}

func main() {
	err := argParse()
	if err != nil {
		log.Fatal(err)
	}
	for {
		log.SetOutput(os.Stdout)
		myactip := getMyIP()
		if checkAllowIP(myactip) == false {
			log.Println("IP not allowed in sucuri firewall, running function to allow")
			addIP()
		}else{
			log.Println("IP is already allowed")
		}
		time.Sleep(120 * time.Second)
	}
}
