package scheduler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/tymofii-polekhin/megos-framework/executor"
	"github.com/tymofii-polekhin/megos-framework/utils"
)

// Scheduler struct that hold all info
type Scheduler struct {
	Master                   string
	FrameworkID              utils.Value
	FrameworkName            string
	FrameworkUser            string
	HeartbeatIntervalSeconds float64
	SchedulerID              string
	MesosStreamID            string
	Transport                *http.Transport
	Client                   *http.Client
	Response                 *http.Response
	RecordIO                 *bufio.Reader
	EventBus                 chan []byte
	Events                   chan utils.Event
	Calls                    chan utils.Call
	TasksLaunched            int
	TasksTotal               int
}

// Subscribe will establish a connection to Mesos master
// and try to register a new framework
func (s *Scheduler) Subscribe(master string, name string, user string) (err error) {

	// fill scheduler structure
	s.Master = master
	s.FrameworkUser = user
	s.FrameworkName = name

	// Create a JSON structure for SUBSCRIBE event
	//log.Println("DEBUG: Create a JSON structure for SUBSCRIBE event")
	call := new(utils.Call)
	call.Type = "SUBSCRIBE"
	call.Subscribe = new(utils.Subscribe)
	call.Subscribe.FrameworkInfo.Name = s.FrameworkName
	call.Subscribe.FrameworkInfo.User = s.FrameworkUser

	url := "http://" + s.Master + "/api/v1/scheduler"

	// Build a JSON string from struct
	//log.Println("DEBUG: Build a JSON string from struct")
	jreq, err := json.Marshal(call)
	if err != nil {
		return err
	}
	//log.Println("DEBUG: string(jreq):", string(jreq))

	// Create POST body from JSON string
	//log.Println("DEBUG: Create POST body from JSON string")
	body := bytes.NewBuffer(jreq)

	// Form an http POST request
	//log.Println("DEBUG: Form an http POST request")
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	// Adding required headers
	req.Header.Set("Host", "localhost:5050")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Connection", "keep-alive")

	// create transport with keep-alive
	//log.Println("DEBUG: create transport with keep-alive")
	s.Transport = &http.Transport{DisableKeepAlives: false}

	// create new http client with keep-alive
	//log.Println("DEBUG: create new http client with keep-alive")
	s.Client = &http.Client{Transport: s.Transport}

	// Do a POST request with SUBSCRIBE message
	//log.Println("DEBUG: Do a POST request with SUBSCRIBE message")
	s.Response, err = s.Client.Do(req)
	if err != nil {
		return err
	}
	//log.Println("DEBUG: Response.ContentLength", s.Response.ContentLength)
	//log.Println("DEBUG: Response.StatusCode", s.Response.StatusCode)
	//log.Println("DEBUG: Response.Status", s.Response.Status)

	s.MesosStreamID = s.Response.Header.Get("Mesos-Stream-Id")
	s.RecordIO = bufio.NewReader(s.Response.Body)
	s.EventBus = make(chan []byte)
	s.Calls = make(chan utils.Call, 10)
	s.Events = make(chan utils.Event, 10)

	// run RecordIOParser in separate goroutine
	//log.Println("DEBUG: run RecordIOParser in separate goroutine")
	go s.RecordIOParser()

	return nil
}

// Teardown all tasks, kill executors and unregister framework
func (s *Scheduler) Teardown() (err error) {

	call := new(utils.Call)
	call.Type = "TEARDOWN"
	call.FrameworkID = &s.FrameworkID

	url := "http://" + s.Master + "/api/v1/scheduler"
	err = call.Send(url, s.MesosStreamID)
	if err != nil {
		return err
	}

	// all clear
	return nil
}

// Shutdown executor and kill all its tasks
func (s *Scheduler) Shutdown(executor *executor.Executor) (err error) {

	call := new(utils.Call)
	call.Type = "SHUTDOWN"
	call.Shutdown = new(utils.Shutdown)
	call.FrameworkID = &s.FrameworkID
	call.Shutdown.ExecutorID = executor.ExecutorID
	call.Shutdown.AgentID = executor.AgentID

	url := "http://" + s.Master + "/api/v1/scheduler"
	err = call.Send(url, s.MesosStreamID)
	if err != nil {
		return err
	}

	// all clear
	return nil
}

// RecordIOParser will parse RecordIO stream from master
// and populate strings to EventBus
func (s *Scheduler) RecordIOParser() (err error) {
	for {

		// reading RecordIO message length
		//log.Println("DEBUG: reading RecordIO message length")
		strRecordIOLength, err := s.RecordIO.ReadString('\n')
		if err != nil {
			return err
		}

		// converting string to integer (dropping '\n' symbol)
		//log.Println("DEBUG: converting string to integer")
		intRecordIOLength, err := strconv.ParseInt(strRecordIOLength[:len(strRecordIOLength)-1], 10, 0)
		if err != nil {
			return err
		}

		// creating buffer for RecordIO message
		//log.Println("DEBUG: creating buffer for RecordIO message")
		bRecordIOMessage := make([]byte, intRecordIOLength)

		// reading message to RecordIO buffer
		//log.Println("DEBUG: reading message to RecordIO buffer")
		iBytesRead, err := s.RecordIO.Read(bRecordIOMessage)
		if err != nil {
			return err
		}

		// if message contains something - send it to Events
		if iBytesRead > 0 {
			var event utils.Event
			//log.Println("DEBUG: unmarshal bRecordIOMessage to utils.Event")
			err := json.Unmarshal(bRecordIOMessage, &event)
			if err != nil {
				log.Fatalln("Failed to unmarshal RecordIO message:")
				log.Fatalln(bRecordIOMessage)
			} else {
				//log.Println("DEBUG: send utils.Event to scheduler.Events channel")
				s.Events <- event
			}
		}

	} // for loop
}

// Decline offer from mesos agent
func (s *Scheduler) Decline(offerIDs []utils.Value, refuseSeconds float64) (err error) {

	log.Println("Declining offers", offerIDs)

	call := new(utils.Call)
	call.Type = "DECLINE"
	call.Decline = new(utils.Decline)
	call.FrameworkID = &s.FrameworkID
	call.Decline.Filters.RefuseSeconds = refuseSeconds
	for _, offerID := range offerIDs {
		call.Decline.OfferIDs = append(call.Decline.OfferIDs, offerID)
	}

	url := "http://" + s.Master + "/api/v1/scheduler"
	err = call.Send(url, s.MesosStreamID)
	if err != nil {
		return err
	}

	// all clear
	return nil
}

// Accept mesos offer and launch task
func (s *Scheduler) Accept(accept *utils.Accept) (err error) {

	log.Println("Accepting offers", accept.OfferIDs)

	call := new(utils.Call)
	call.Type = "ACCEPT"
	call.FrameworkID = &s.FrameworkID
	call.Accept = accept

	url := "http://" + s.Master + "/api/v1/scheduler"
	err = call.Send(url, s.MesosStreamID)
	if err != nil {
		return err
	}

	// all clear
	return nil
}

// Acknowledge mesos master event
func (s *Scheduler) Acknowledge(event utils.Event) (err error) {

	call := new(utils.Call)
	call.Type = "ACKNOWLEDGE"
	call.Acknowledge = new(utils.Acknowledge)
	call.FrameworkID = &s.FrameworkID
	call.Acknowledge.AgentID = event.Update.Status.AgentID
	call.Acknowledge.TaskID = event.Update.Status.TaskID
	call.Acknowledge.UUID = event.Update.Status.UUID

	url := "http://" + s.Master + "/api/v1/scheduler"
	err = call.Send(url, s.MesosStreamID)
	if err != nil {
		return err
	}

	// all clear
	return nil
}
