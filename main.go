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
	"log"

	"github.com/tymofii-polekhin/megos-framework/scheduler"
)

//func Scheduler()

//func Register()
//func Reregister()
//func Deregister()
//func Message()
//func Heartbeat()

// Scheduler info

// EventHandler will distribute events to corresponding func
// func EventHandler(s *utils.Scheduler) {
// 	for {
// 		e := <-s.EventBus
// 		//log.Println("EventHandler received an event:")
// 		//log.Println(string(e))
// 		var j utils.SchedulerEvent
// 		err := json.Unmarshal(e, &j)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		//log.Println(j.Type)
// 		switch {
//
// 		case j.Type == "SUBSCRIBED":
// 			log.Println("Subscribed to master at", s.MesosMaster,
// 				"with framework ID", j.Subscribed.FrameworkID.Value,
// 				"and heartbeat interval", j.Subscribed.HeartbeatIntervalSeconds, "seconds")
//
// 		case j.Type == "OFFERS":
// 			// Ranging through offers from multiple agents
// 			for _, o := range j.Offers.Offers {
// 				log.Println("\tReceived offer", o.ID.Value)
// 				log.Println("\tAgent", o.AgentID.Value, "on host", o.Hostname, "with:")
// 				s.Offers = append(s.Offers, o)
// 				// Ranging throgh multiple resources from single agent
// 				for _, r := range o.Resources {
// 					if r.Type == "SCALAR" {
// 						log.Println("\t\t", r.Name, r.Scalar.Value)
// 					}
// 					if r.Type == "RANGES" {
// 						for _, p := range r.Ranges.Range {
// 							log.Println("\t\t", r.Name, p.Begin, "-", p.End)
// 						}
// 					}
// 				}
// 			}
// 			SchedulerOfferProcess(s)
//
// 		case j.Type == "HEARTBEAT":
// 			log.Println("Heartbeat received")
//
// 		}
// 	}
// }
//
// // SchedulerOfferProcess process offers from cluster
// func SchedulerOfferProcess(s *utils.Scheduler) {
// 	log.Println("SchedulerOfferProcess: received offer")
//
// 	for _, o := range s.Offers {
// 		log.Println("Processing DECLINE for offer", o.ID.Value)
// 		SchedulerOfferDecline(o.ID.Value, o.FrameworkID.Value, s.MesosStreamID)
// 		// need to delete offer from Scheduler.Offers
// 		s.Offers = s.Offers[1:]
// 	}
//
// }
//
// // SchedulerOfferDecline declining one offer
// func SchedulerOfferDecline(offerID string, frameworkID string, mesosSreamID string) {
//
// 	var value utils.Value
// 	value.Value = offerID
//
// 	var message utils.DeclineOffer
// 	message.Type = "DECLINE"
// 	message.FrameworkID.Value = frameworkID
// 	message.Decline.Filters.RefuseSeconds = 5.0
// 	message.Decline.OfferIDs = append(message.Decline.OfferIDs, value)
//
// 	url := "http://localhost:5050/api/v1/scheduler"
// 	//log.Println("Subscribe: url:", url)
//
// 	jreq, err := json.Marshal(message)
// 	body := bytes.NewBuffer(jreq)
// 	//log.Println("Subscribe: body:", body)
//
// 	// Forming http POST request to subscribe
// 	req, err := http.NewRequest("POST", url, body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Adding headers
// 	req.Header.Set("Host", "localhost:5050")
// 	req.Header.Set("Content-type", "application/json")
// 	req.Header.Set("Mesos-Stream-Id", mesosSreamID)
//
// 	var client = &http.Client{}
// 	responce, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer responce.Body.Close()
//
// 	log.Println("Responce from Master:", responce.Status)
//
// }
//
// // Subscribe func
// func Subscribe(s *utils.Scheduler) {
//
// 	var message utils.SubscribeMessage
// 	message.Type = "SUBSCRIBE"
// 	message.Subscribe.FrameworkInfo.Name = "Go Sample Framework"
// 	message.Subscribe.FrameworkInfo.User = "root"
//
// 	url := "http://" + s.MesosMaster + "/api/v1/scheduler"
// 	//log.Println("Subscribe: url:", url)
//
// 	jreq, err := json.Marshal(message)
// 	body := bytes.NewBuffer(jreq)
// 	//log.Println("Subscribe: body:", body)
//
// 	// Forming http POST request to subscribe
// 	req, err := http.NewRequest("POST", url, body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Adding headers
// 	req.Header.Set("Host", "localhost:5050")
// 	req.Header.Set("Content-type", "application/json")
// 	req.Header.Set("Connection", "keep-alive")
//
// 	// Do a POST request with SUBSCRIBE message
// 	s.Response, err = s.MainConn.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Reading never-ending stream from Mesos master
// 	s.MesosStreamID = s.Response.Header.Get("Mesos-Stream-Id")
// 	s.RecordIO = bufio.NewReader(s.Response.Body)
// }
//
// // RecordIOParse function populates clean messages to EventBus
// func RecordIOParse(s *utils.Scheduler) {
// 	for {
// 		line, err := s.RecordIO.ReadString('\n')
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		i, err := strconv.ParseInt(line[:len(line)-1], 10, 0)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		msg := make([]byte, i)
// 		n, err := s.RecordIO.Read(msg)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if n > 0 {
// 			s.EventBus <- msg
// 		}
// 	} // for loop
// }

func main() {

	// create new Scheduler
	s := new(scheduler.Scheduler)

	// subscribe
	serr := s.Subscribe("localhost:5050", "New Go Scheduler", "root")
	if serr != nil {
		log.Fatalln(serr)
	}

	// RecordIO parser
	go s.RecordIOParser()

	// EventHandler
	go s.EventHandler()

	//forever loop
	for {
	}
}
