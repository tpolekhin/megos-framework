package scheduler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tymofii-polekhin/megos-framework/utils"
)

// Scheduler struct that hold all info
type Scheduler struct {
	Master                   string
	FrameworkID              string
	FrameworkName            string
	FrameworkUser            string
	HeartbeatIntervalSeconds float64
	SchedulerID              string
	MesosStreamID            string
	MainTransport            *http.Transport
	MainConn                 *http.Client
	Response                 *http.Response
	RecordIO                 *bufio.Reader
	EventBus                 chan []byte
	Offers                   []utils.Offer
}

// Subscribe will establish a connection to Mesos master
// and try to register a new framework
func (s *Scheduler) Subscribe(master string, name string, user string) (err error) {

	// Create a JSON structure for SUBSCRIBE event
	var message utils.SubscribeMessage
	message.Type = "SUBSCRIBE"
	message.Subscribe.FrameworkInfo.Name = name
	message.Subscribe.FrameworkInfo.User = user

	url := "http://" + master + "/api/v1/scheduler"

	// Build a JSON string from struct
	jreq, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Create POST body from JSON string
	body := bytes.NewBuffer(jreq)

	// Form an http POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	// Adding required headers
	req.Header.Set("Host", "localhost:5050")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Connection", "keep-alive")

	// create transport with keep-alive
	s.MainTransport = &http.Transport{DisableKeepAlives: false}

	// create new http client with keep-alive
	s.MainConn = &http.Client{Transport: s.MainTransport}

	// Do a POST request with SUBSCRIBE message
	s.Response, err = s.MainConn.Do(req)
	if err != nil {
		return err
	}

	s.MesosStreamID = s.Response.Header.Get("Mesos-Stream-Id")
	s.RecordIO = bufio.NewReader(s.Response.Body)
	s.Master = master
	s.FrameworkUser = user
	s.FrameworkName = name
	s.EventBus = make(chan []byte)

	return nil
}

// RecordIOParser will parse RecordIO stream from master
// and populate strings to EventBus
func (s *Scheduler) RecordIOParser() (err error) {
	for {

		// reading RecordIO message length
		strRecordIOLength, err := s.RecordIO.ReadString('\n')
		if err != nil {
			return err
		}

		// converting string to integer (dropping '\n' symbol)
		intRecordIOLength, err := strconv.ParseInt(strRecordIOLength[:len(strRecordIOLength)-1], 10, 0)
		if err != nil {
			return err
		}

		// creating buffer for RecordIO message
		bRecordIOMessage := make([]byte, intRecordIOLength)

		// reading message to RecordIO buffer
		iBytesRead, err := s.RecordIO.Read(bRecordIOMessage)
		if err != nil {
			return err
		}

		// if message contains something - send it to EventBus
		if iBytesRead > 0 {
			s.EventBus <- bRecordIOMessage
		}

	} // for loop
}

// EventHandler will read messages from EventBus and
// handle for further processing
func (s *Scheduler) EventHandler() (err error) {
	for {

		// read next event from EventBus
		e := <-s.EventBus

		//log.Println("EventHandler received an event:")
		//log.Println(string(e))

		// create SchedulerEvent structure
		var j utils.SchedulerEvent

		// try to parse event into JSON structure
		err := json.Unmarshal(e, &j)
		if err != nil {
			return err
		}

		// choosing handler based on event type
		switch {

		case j.Type == "SUBSCRIBED": // master has accepted registration
			s.FrameworkID = j.Subscribed.FrameworkID.Value
			s.HeartbeatIntervalSeconds = j.Subscribed.HeartbeatIntervalSeconds
			log.Println("Subscribed to master at", s.Master,
				"with framework ID", j.Subscribed.FrameworkID.Value,
				"and heartbeat interval", j.Subscribed.HeartbeatIntervalSeconds, "seconds")

		case j.Type == "OFFERS": // incomig resource offers
			// Ranging through offers from multiple agents
			for _, o := range j.Offers.Offers {
				log.Println("\tReceived offer", o.ID.Value)
				log.Println("\tAgent", o.AgentID.Value, "on host", o.Hostname, "with:")
				s.Offers = append(s.Offers, o)
				// Ranging throgh multiple resources from single agent
				for _, r := range o.Resources {
					if r.Type == "SCALAR" {
						log.Println("\t\t", r.Name, r.Scalar.Value)
					}
					if r.Type == "RANGES" {
						for _, p := range r.Ranges.Range {
							log.Println("\t\t", r.Name, p.Begin, "-", p.End)
						}
					}
				}
			}
			s.SchedulerOfferProcess()

		case j.Type == "HEARTBEAT": // master heartbeat event
			log.Println("Heartbeat received")

		}
	}
}

// SchedulerOfferProcess process offers from cluster
func (s *Scheduler) SchedulerOfferProcess() (err error) {
	log.Println("SchedulerOfferProcess: received offer")

	for _, o := range s.Offers {
		log.Println("Processing DECLINE for offer", o.ID.Value)
		e := s.DeclineOffer(o.ID.Value)
		if err != nil {
			return e
		}
		// need to delete offer from Scheduler.Offers
		s.Offers = s.Offers[1:]
	}

	return nil
}

// DeclineOffer declining one offer
func (s *Scheduler) DeclineOffer(offerID string) (err error) {

	var value utils.Value
	value.Value = offerID

	var message utils.DeclineOffer
	message.Type = "DECLINE"
	message.FrameworkID.Value = s.FrameworkID
	message.Decline.Filters.RefuseSeconds = 5.0
	message.Decline.OfferIDs = append(message.Decline.OfferIDs, value)

	url := "http://" + s.Master + "/api/v1/scheduler"

	jreq, e := json.Marshal(message)
	if e != nil {
		return e
	}
	body := bytes.NewBuffer(jreq)

	// Forming http POST request to subscribe
	req, e := http.NewRequest("POST", url, body)
	if e != nil {
		return e
	}

	// Adding headers
	req.Header.Set("Host", "localhost:5050")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Mesos-Stream-Id", s.MesosStreamID)

	client := new(http.Client)
	responce, e := client.Do(req)
	if e != nil {
		return e
	}
	defer responce.Body.Close()

	log.Println("Responce from Master:", responce.Status)
	if responce.StatusCode != http.StatusAccepted {
		return fmt.Errorf(responce.Status)
	}
	return nil
}
