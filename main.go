package main

import (
	"github.com/go-co-op/gocron"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type myIpResponse struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
}

type dyndnsResponse struct {
	Success bool
	Message string
}

func getMyIpResponse(url string) (myIpResponse, error) {
	response := myIpResponse{}

	log.Println("Getting ", url)
	res, err := http.Get(url)
	if err != nil {
		return response, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func getDyndnsReponse(url string) (dyndnsResponse, error) {
	response := dyndnsResponse{}

	res, err := http.Get(url)
	if err != nil {
		return response, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func main() {
	url1 := "https://api4.my-ip.io/ip.json"
	secret := os.Getenv("DYNDNS_SECRET")
	subdomain := os.Getenv("DYNDNS_SUBDOMAIN")

	s := gocron.NewScheduler(time.UTC)

	s.Every("10m").Do(func() {
		fmt.Println("Starting")

		res1, err := getMyIpResponse(url1)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("IP response was ", res1)
		if res1.Success {
			url2 := fmt.Sprintf("https://ns.kleinklein.net/update?secret=%s&domain=%s&addr=%s", secret, subdomain, res1.IP)

			res2, err := getDyndnsReponse(url2)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("DNS response was ", res2)
		} else {
			fmt.Println("Not sending DNS request due to IP failure")
		}
	})

	s.StartBlocking()
}
