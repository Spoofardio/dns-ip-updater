package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	namecheap "github.com/billputer/go-namecheap"
	"github.com/jasonlvhit/gocron"
	"github.com/sfreiberg/gotwilio"
)

// IPChangeRequest object for POST on /ip
type IPChangeRequest struct {
	OldIP string `json:"old_ip"`
	NewIP string `json:"new_ip"`
	Key   string `json:"key"`
}

var namecheapClient *namecheap.Client
var twilio *gotwilio.Twilio

var serverKey string
var serverUsername string
var serverPassword string
var currentIP string

var ipv4Regex *regexp.Regexp

func main() {
	serverKey = os.Getenv("KEY")
	serverUsername = os.Getenv("USERNAME")
	serverPassword = os.Getenv("PASSWORD")

	ipv4Regex, _ = regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	twilio = gotwilio.NewTwilioClient(twilioSid, twilioToken)
	currentIP = getIP()
	log.Println("DNS IP Updater Client Started: ", time.Now().Format("2006-Jan-02"))

	gocron.Every(1).Minutes().Do(checkIP)
	<-gocron.Start()
}

func checkIP() {
	ip := getIP()
	if ip != currentIP {
		message := "Host's IP address has changed to " + string(ip)
		twilio.SendSMS(twilioFrom, twilioTo, message, "", "")
		log.Println(message+" at ", time.Now().Format("2006-Jan-02"))
		sendDNSUpdateRequest(ip)
		currentIP = ip
	}
}

func getIP() string {
	res, _ := http.Get("https://api.ipify.org")
	ip, _ := ioutil.ReadAll(res.Body)
	isValidIP := ipv4Regex.MatchString(string(ip))
	if !isValidIP {
		log.Printf("Received invalid IP: %s", string(ip))
		return currentIP
	}
	return string(ip)
}

func sendDNSUpdateRequest(ip string) {
	request := &IPChangeRequest{
		OldIP: currentIP,
		NewIP: ip,
		Key:   serverKey,
	}
	body, err := json.Marshal(request)
	if err != nil {
		log.Panic(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(serverUsername, serverPassword)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("The DNS Records were succesfully Updated")
	} else {
		log.Panicf("There was an error with updating the DNS Records: " + resp.Status)
	}
}
