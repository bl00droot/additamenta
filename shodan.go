package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"net"
	"time"
)

type shodaneiro struct {
	More   string `json:"more"`
	Domain string `json:"domain"`
	Tags   string `json:"tags"`
	Data   struct {
		Subdomain string `json:"subdomain"`
		Type      string `json:"type"`
		Value     string `json:"value"`
		Last      string `json:"last_seen"`
	}
	Subdomains []string `json:"subdomains"`
}


func checkConn(ip string, portas string){
	address := net.JoinHostPort(ip, portas)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		//fmt.Printf("DOMAIN:%s:%s Fechado\n", ip, portas)
	} else {
		if conn != nil {
			portsOpen = append(portsOpen, portas)
			//fmt.Printf("DOMAIN:%s:%s Aberto\n", ip, portas)
			_ = conn.Close()
			
		}
	}
}



var (
	exit   = os.Exit // typedef like
	apikey = flag.String("apikey", "", "apikey: YourApiKey - obrigatorio")
	domain = flag.String("domain", "", "domain for subdomain mapping - obrigatorio")
	defPorts = []int{66, 80, 81, 443, 10443, 9443, 445, 457, 1080, 1100, 1241, 1352, 1433, 1434, 1521, 1944, 2301, 3128, 3306, 4000, 4001, 4002, 4100, 5000, 5432, 5800, 5801, 5802, 6346, 6347, 7001, 7002, 8080, 8081, 8181, 8888, 30821}
	portsOpen  []string
)

func main() {

	flag.Parse()

	client := &http.Client{}

	if *apikey == "" {
		flag.Usage()
		exit(1)
	}
	if *domain == "" {
		flag.Usage()
		exit(1)
	}
	req, err := http.NewRequest("GET", "https://api.shodan.io/dns/domain/"+*domain+"?key="+*apikey, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/6.0 (X11; Linux x86) AppleWebKit/539.666 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/539.666")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erro na response", err)
	}
	//var portsOpen []string
	var data shodaneiro
	err = json.Unmarshal([]byte(body), &data)
	//fmt.Printf("%s\n", data.Subdomains)
	for i := 0; i < len(data.Subdomains); i++ {
		//fmt.Println(data.Matches[i].IPStr, ":", data.Matches[i].Port)
		//fmt.Printf("%s.%s\n", data.Subdomains[i], *domain)
		subdomaincheck := fmt.Sprintf("%s.%s", data.Subdomains[i], *domain)
		for p := 0; p < len(defPorts); p++ {
			porta := strconv.Itoa(defPorts[p])
			checkConn(subdomaincheck, porta)
			
		}

		fmt.Printf("DOMAIN:%s | %s \n", subdomaincheck, portsOpen)
		portsOpen = nil
		//var portsOpen []string = nil
		//portsOpen = append(portsOpen, portas)
		//clearSlice()
		//fmt.Println("STATUS CODE HTTPS: ", resphttps.StatusCode, http.StatusText(resphttps.StatusCode))

	}

	exit(0)

}