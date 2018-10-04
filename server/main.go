package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	namecheap "github.com/billputer/go-namecheap"
	"github.com/gorilla/mux"
)

// IPChangeRequest object for POST on /ip
type IPChangeRequest struct {
	OldIP string `json:"old_ip"`
	NewIP string `json:"new_ip"`
	Key   string `json:"key"`
}

var namecheapClient *namecheap.Client
var connectionKey string

func main() {
	connectionKey = os.Getenv("KEY")

	// namecheap setup
	namecheapClient = namecheap.NewClient(namecheapAPIUser, namecheapAPIToken, namecheapUserName)
	log.Println("DNS IP Updater Server Started")
	log.Println("Now serving requests on \"/\" and \"/ip\"")

	// http setup
	router := mux.NewRouter()
	router.HandleFunc("/", healthCheck).Methods("GET")
	router.HandleFunc("/ip", ipChange).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// healthCheck Handler for /
// a simple health check that returns OK 200
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

// ipChange Handler for /ip
// updated namecheap DNS hosts to be the new IP.
func ipChange(w http.ResponseWriter, r *http.Request) {
	log.Println("Ip Change request")
	var request IPChangeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil { // 400
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	if request.Key != connectionKey { // 401
		log.Println("Key: " + request.Key + " is not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not Authorized"))
		return
	}

	hostAddressChanges := "Ip Change from: " + request.OldIP + " --> " + request.NewIP + "\nUpdated Hosts: "
	dns, err := namecheapClient.DomainsDNSGetHosts(namecheapSLD, namecheapTLD)

	if err != nil { // 500
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	hosts := dns.Hosts
	for i := 0; i < len(hosts); i++ {
		host := &hosts[i]
		if host.Address == request.OldIP {
			host.Address = request.NewIP
			hostAddressChanges += host.Name + " "
		}
	}

	res, _ := namecheapClient.DomainDNSSetHosts(namecheapSLD, namecheapTLD, hosts)

	if res.IsSuccess != true { // 502
		log.Println("Namecheap error while updating the DNS records")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Namecheap error while updating the DNS records"))
		return
	}

	log.Printf(hostAddressChanges) // 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success\n" + hostAddressChanges))
}
