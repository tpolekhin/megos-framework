package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Calls: Scheduler -> Master
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

// Call message from scheduler/executor to mesos
type Call struct {
	Type        string       `json:"type"`
	FrameworkID *Value       `json:"framework_id,omitempty"`
	Subscribe   *Subscribe   `json:"subscribe,omitempty"`
	Accept      *Accept      `json:"accept,omitempty"`
	Decline     *Decline     `json:"decline,omitempty"`
	Revive      *Revive      `json:"revive,omitempty"`
	Kill        *Kill        `json:"kill,omitempty"`
	Shutdown    *Shutdown    `json:"shutdown,omitempty"`
	Acknowledge *Acknowledge `json:"acknowledge,omitempty"`
	Reconcile   *Reconcile   `json:"reconcile,omitempty"`
	Message     *Message     `json:"message,omitempty"`
	Request     *Request     `json:"request,omitempty"`
}

// Revive comment
type Revive struct{}

// Kill comment
type Kill struct{}

// Reconcile comment
type Reconcile struct{}

// Message comment
type Message struct{}

// Request comment
type Request struct{}

// Send call to mesos master/agent
func (c *Call) Send(endpoint string, streamID string) (err error) {

	// Create a JSON string from Call structure
	message, e := json.Marshal(c)
	if e != nil {
		return e
	}

	// Create a new http request byte buffer
	body := bytes.NewBuffer(message)

	// Forming http POST request to subscribe
	request, e := http.NewRequest("POST", endpoint, body)
	if e != nil {
		return e
	}

	// Adding required headers
	request.Header.Set("Host", "localhost:5050")
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Mesos-Stream-Id", streamID)

	// Create new http client
	client := new(http.Client)

	// do http POST to endpoint
	res, e := client.Do(request)
	if e != nil {
		return e
	}

	// raise an error if mesos didn't replied 202 Accepted
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf(res.Status, string(message))
	}

	// all clear
	return nil
}

// Subscribe comment
type Subscribe struct {
	FrameworkInfo FrameworkInfo `json:"framework_info"`
}

// Accept comment
type Accept struct {
	OfferIDs   []Value     `json:"offer_ids"`
	Operations []Operation `json:"operations"`
	Filters    Filters     `json:"filters"`
}

// Decline comment
type Decline struct {
	OfferIDs []Value `json:"offer_ids"`
	Filters  Filters `json:"filters"`
}

// Shutdown comment
type Shutdown struct {
	ExecutorID Value `json:"executor_id"`
	AgentID    Value `json:"agent_id"`
}

// Acknowledge comment
type Acknowledge struct {
	AgentID Value  `json:"agent_id"`
	TaskID  Value  `json:"task_id"`
	UUID    string `json:"uuid"`
}
