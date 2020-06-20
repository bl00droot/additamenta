package main

import (
	"fmt"
	"bufio"
	"os"
	"flag"
	"log"
	"net/http"
	"io/ioutil"	
	"encoding/json"
	"net"
	"time"
	"strconv"
)

type shodaneiro struct {
	IP          string        `json:"ip"`
	Ports        []int       `json:"ports"`
}

var (
	exit    = os.Exit 
	filePtr = flag.String("file", "", "ip list - obrigatorio")
	apikey  = flag.String("apikey", "", "apikey: YourApiKey - obrigatorio")
)



func readFile() {

	ipFile, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal("Erro na leitura", err)	
	}
	defer ipFile.Close()
	scanner := bufio.NewScanner(ipFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		ipStr = append(ipsStr, scanner.Text())
	} 
}



func checkConn(ip string, portas string){
	address := net.JoinHostPort(ip, portas)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		fmt.Printf("IP:%s:%s Fechado\n", ip, portas)
	} else {
		if conn != nil {
			fmt.Printf("IP:%s:%s Aberto\n", ip, portas)
			_ = conn.Close()
		}
	}
}


func main() {

	flag.Parse()

	var ipsStr []string
	client := &http.Client{}

	if *filePtr == "" {
		flag.Usage()
		exit(1)
	}
/***
    if *filePtr != "" {
		ipFile, err := os.Open(*filePtr)
		if err != nil {
			log.Fatal("Erro na leitura", err)
		}
		defer ipFile.Close()
		scanner := bufio.NewScanner(ipFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
		  ipsStr = append(ipsStr, scanner.Text())
		}
	  }
***/
	

	if *apikey == "" {
		flag.Usage()
		exit(1)
	}

	readFile()
	
	for i := range ipsStr {
		ipsEnt :=  ipsStr[i]

		req, err := http.NewRequest("GET", "https://api.shodan.io/shodan/host/"+ipsEnt+"?key="+*apikey+"&history=yes", nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Erro na response", err)
		}
		var data shodaneiro
		err = json.Unmarshal([]byte(body), &data)
		//fmt.Printf("IP:%s:%d \n", ipsEnt, data.Ports)
		//results := make()
		for p := 0; p < len(data.Ports); p++ {
			portaTo := strconv.Itoa(data.Ports[p])
			checkConn(ipsEnt, portaTo)

		}

		

	}
	exit(0)
}
  