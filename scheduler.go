// GOlang implementation of mesos framework HTTP interface
//
// Messaging: Scheduler -> Master
// SUBSCRIBE
// TEARDOWN
// ACCEPT
// DECLINE
// REVIVE
// KILL
// SHUTDOWN
// ACKNOWLEDGE
// RECONCILE
// MESSAGE
// REQUEST
//
// Events: Master -> Scheduler
// SUBSCRIBED
// OFFERS
// RESCIND
// UPDATE
// MESSAGE
// FAILURE
// ERROR
// HEARTBEAT

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tpolekhin/mesos-1.0.0-http-scheduler/utils"
)

//func Scheduler()

//func Register()
//func Reregister()
//func Deregister()
//func Message()
//func Heartbeat()

// Scheduler info
var Scheduler utils.Scheduler

// EventHandler will distribute events to corresponding func
func EventHandler(event chan []byte) {
	for {
		e := <-event
		var j utils.SchedulerEvent
		err := json.Unmarshal(e, &j)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(j.Type)
		switch {
		case j.Type == "SUBSCRIBED":
			Scheduler.SchedulerID = j.Subscribed.FrameworkID.Value
			fmt.Println(j.Subscribed.FrameworkID.Value)
			fmt.Println(j.Subscribed.HeartbeatIntervalSeconds)
		case j.Type == "OFFERS":
			for _, o := range j.Offers.Offers {
				fmt.Println(o.AgentID.Value)
				fmt.Println(o.FrameworkID.Value)
				fmt.Println(o.Hostname)
				for _, r := range o.Resources {
					fmt.Println(r.Name)
				}
			}
		case j.Type == "HEARTBEAT":
			fmt.Println("Heartbeat received from master.")
		}
	}
}

// SchedulerOfferProcess process offers from cluster
func SchedulerOfferProcess(o utils.SchedulerEvent) {
	fmt.Println("SchedulerOfferProcess: received offer")

	for _, o := range o.Offers.Offers {
		fmt.Println("Processing DECLINE for offer ", o.ID.Value)
		SchedulerOfferDecline(o.ID.Value)
	}

}

// SchedulerOfferDecline declining one offer
func SchedulerOfferDecline(offerID string) {

	body := bufio.Reader

	req, err := http.NewRequest("POST", "http://localhost:5050/api/v1/scheduler")
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Host", "localhost:5050")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Mesos-Stream-Id", Scheduler.MesosStreamID)

}

func main() {

	events := make(chan []byte)
	go EventHandler(events)

	url := "http://localhost:5050/api/v1/scheduler"
	//fmt.Println("URL:> ", url)

	// Create a custom transport with keep-alive enabled
	tr := &http.Transport{DisableKeepAlives: false}

	// Create client using custom keep-alive transport
	client := &http.Client{Transport: tr}

	// form a post message
	var jsonStr = []byte(`{"type":"SUBSCRIBE","subscribe":{"framework_info":{"user":"root","name":"GOlang HTTP Framework"}}}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	// adding headers
	req.Header.Set("Host", "localhost:5050")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Connection", "keep-alive")

	//fmt.Println("req:> ", req)

	// Do a POST request with SUBSCRIBE message and wait for responce
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)
	Scheduler.MesosStreamID = resp.Header.Get("Mesos-Stream-Id")
	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("RecordIO Message Size:>", line[:len(line)-1])

		i, err := strconv.ParseInt(line[:len(line)-1], 10, 0)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("strconv.ParseInt:>", i)

		msg := make([]byte, i)
		n, err := reader.Read(msg)
		if err != nil {
			fmt.Println(err)
		}
		if n > 0 {
			events <- msg
		}

		//fmt.Println("Bytes read:>", n)
		//message := string(msg)
		//fmt.Println("String read:>", message)
		//fmt.Println("Sending message to channel...")
		//events <- message

	} // for loop

}
