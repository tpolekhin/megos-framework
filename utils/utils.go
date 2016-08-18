package utils

import (
	"bufio"
	"fmt"
	"net/http"
)

// Scheduler type that holds schduler specific information
type Scheduler struct {
	MesosMaster   string
	SchedulerID   string
	MesosStreamID string
	MainConn      *http.Client
	RecordIO      *bufio.Reader
	EventBus      chan []byte
}

// SubscribeMessage comment
type SubscribeMessage struct {
	Type      string    `json:"type"`
	Subscribe subscribe `json:"subscribe"`
}

type subscribe struct {
	FrameworkInfo frameworkInfo `json:"framework_info"`
}

type frameworkInfo struct {
	User string `json:"user"`
	Name string `json:"name"`
}

// SchedulerEvent comment
type SchedulerEvent struct {
	Type       string     `json:"type"`
	Subscribed subscribed `json:"subscribed,omitempty"`
	Offers     offers     `json:"offers,omitempty"`
}

// subscribed comment
type subscribed struct {
	HeartbeatIntervalSeconds float64     `json:"heartbeat_interval_seconds"`
	FrameworkID              frameworkID `json:"framework_id"`
}

// framework_id comment
type frameworkID struct {
	Value string `json:"value"`
}

type offers struct {
	Offers []offer `json:"offers"`
}

type offer struct {
	AgentID     agentID     `json:"agent_id"`
	FrameworkID frameworkID `json:"framework_id"`
	Hostname    string      `json:"hostname"`
	ID          id          `json:"id"`
	Resources   []resource  `json:"resources,omitempty"`
	URL         url         `json:"url,omitempty"`
}

type agentID struct {
	Value string `json:"value"`
}

type id struct {
	Value string `json:"value"`
}

type url struct {
	Address address `json:"address"`
	Path    string  `json:"path"`
	Scheme  string  `json:"scheme"`
}

type address struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

type resource struct {
	Name   string  `json:"name"`
	Role   string  `json:"role"`
	Type   string  `json:"type"`
	Scalar scalar  `json:"scalar,omitempty"`
	Ranges rangess `json:"ranges,omitempty"`
}

type scalar struct {
	Value float32 `json:"value"`
}

type rangess struct {
	Range []ranges `json:"ranges"`
}

type ranges struct {
	Begin float64 `json:"begin"`
	End   float64 `json:"end"`
}

// DeclineOffer struct to send to master when declining offer
type DeclineOffer struct {
	FrameworkID frameworkID `json:"framework_id"`
	Type        string      `json:"type"`
	Decline     decline     `json:"decline"`
}

type decline struct {
	OfferIDs []value `json:"offer_ids"`
	Filters  filters `json:"filters"`
}

type filters struct {
	RefuseSeconds float64 `json:"refuse_seconds"`
}

type value struct {
	Value string `json:"value"`
}

func jsonloop(i interface{}) {
	m := i.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Print(k, "(string):", vv, "\n")
		case int:
			fmt.Print(k, "(int):", vv, "\n")
		case float64:
			fmt.Print(k, "(float64):", vv, "\n")
		case []interface{}:
			fmt.Print(k, "(array):\n")
			for _, u := range vv {
				jsonloop(u)
			}
		case interface{}:
			fmt.Print(k, "(object):\n")
			jsonloop(v)
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}

func jsonpprint(i interface{}, indent int) {
	tab := ""
	for n := 0; n < indent; n++ {
		tab += "  "
	}
	m := i.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Print(tab, k, ": ", vv, "\n")
		case int:
			fmt.Print(tab, k, ": ", vv, "\n")
		case float64:
			fmt.Print(tab, k, ": ", vv, "\n")
		case []interface{}:
			fmt.Print(tab, k, ":[\n")
			for _, u := range vv {
				jsonpprint(u, indent+1)
				fmt.Print(tab, "]\n")
			}
		case interface{}:
			fmt.Print(tab, k, ":{\n")
			jsonpprint(v, indent+1)
			fmt.Print(tab, "}\n")
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
}
