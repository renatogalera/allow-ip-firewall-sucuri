package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var API_KEY string
var API_SECRET string
var NEWIPADDR string

func argParse() error {

	// Checa se arquivo config.env existe
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Coletando variáveis
	API_KEY = os.Getenv("API_KEY")
	if API_KEY == "" {
		msg := fmt.Sprintf("É necessário configurar variável API_KEY")
		return errors.New(msg)
	}
	API_SECRET = os.Getenv("API_SECRET")
	if API_SECRET == "" {
		msg := fmt.Sprintf("É necessário configurar variável API_SECRET")
		return errors.New(msg)
	}
	NEWIPADDR = os.Getenv("NEWIPADDR")
	if NEWIPADDR == "" {
		msg := fmt.Sprintf("É necessário configurar variável NEWIPADDR")
		return errors.New(msg)
	}
	// Criando arquivo novoip caso não exista
	if _, err := os.Stat(NEWIPADDR); os.IsNotExist(err) {
		os.OpenFile(NEWIPADDR, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	}

	return nil
}

func checkIP() {
	log.Printf("Checando IP...\n")
	// Pega ip ipv4 lá de baixo
	IPV4 := getMyIP(4)
	// le o IPnovo.txt
	OLDIPTMP, _ := ioutil.ReadFile(NEWIPADDR)
	OLDIP := string(OLDIPTMP)

	if OLDIP != IPV4 {
		log.Printf("IP Mudou! Alterando: %s -> %s", OLDIP, IPV4)
		// Salva arquivo com novo IP!

		//Deletando antigo IP
		var err = os.Remove(NEWIPADDR)
		if err != nil {
			panic(err)
		}

		//Criando arquivo
		saveip, err := os.OpenFile(NEWIPADDR, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		defer saveip.Close()
		if _, err = saveip.WriteString(IPV4); err != nil {
			panic(err)
		}
		addIP(IPV4, OLDIP)

	}
}

func addIP(IPV4, OLDIP string) string {
	fmt.Println("O IP mudou, alterando")
	PAGINA := "https://waf.sucuri.net/api?v2&k=%s&s=%s&a=%s&ip=%s"
	RMIP := fmt.Sprintf(PAGINA, API_KEY, API_SECRET, "delete_whitelist_ip", OLDIP)
	ALLOWIP := fmt.Sprintf(PAGINA, API_KEY, API_SECRET, "whitelist_ip", IPV4)

	client := &http.Client{}

	if OLDIP != "" {
		req, _ := http.NewRequest("GET", RMIP, nil)
		req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
		resp, _ := client.Do(req)
		fmt.Println(resp)
		time.Sleep(2 * time.Second)
	}

	req2, _ := http.NewRequest("GET", ALLOWIP, nil)
	req2.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	resp2, _ := client.Do(req2)
	fmt.Println(resp2)

	return ""

}

func getMyIP(protocol int) string {
	var target string
	if protocol == 4 {
		target = "http://ifconfig.me/ip"

		//} else if protocol == 6 {
		//	target = "http://ifconfig.me/ip"
		//
	} else {
		return ""

	}
	resp, err := http.Get(target)

	if err == nil {
		contents, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			defer resp.Body.Close()
			return strings.TrimSpace(string(contents))

		}

	}
	return ""
}

func main() {

	err := argParse()
	if err != nil {
		log.Fatal(err)
	}

	checkIP()

}
