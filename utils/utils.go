package utils

import "fmt"

// FrameworkInfo comment
type FrameworkInfo struct {
	User string `json:"user"`
	Name string `json:"name"`
}

// Offer struct
type Offer struct {
	AgentID     Value      `json:"agent_id"`
	FrameworkID Value      `json:"framework_id"`
	Hostname    string     `json:"hostname"`
	ID          Value      `json:"id"`
	Resources   []Resource `json:"resources"`
	URL         *URL       `json:"url,omitempty"`
}

// URL comment
type URL struct {
	Address Address `json:"address"`
	Path    string  `json:"path"`
	Scheme  string  `json:"scheme"`
}

//Address comment
type Address struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

// Resource comment
type Resource struct {
	Name   string  `json:"name"`
	Role   string  `json:"role"`
	Type   string  `json:"type"`
	Scalar *Scalar `json:"scalar,omitempty"`
	Ranges *Ranges `json:"ranges,omitempty"`
}

// Scalar comment
type Scalar struct {
	Value float64 `json:"value"`
}

// Ranges comment
type Ranges struct {
	Range []Range `json:"range"`
}

// Range comment
type Range struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}

// Filters comment
type Filters struct {
	RefuseSeconds float64 `json:"refuse_seconds"`
}

// Value struct
type Value struct {
	Value string `json:"value"`
}

// Operation comment
type Operation struct {
	Type   string `json:"type"`
	Launch Launch `json:"launch"`
}

// Launch comment
type Launch struct {
	TaskInfos []TaskInfo `json:"task_infos"`
}

// Executor comment
type Executor struct {
	ExecutorID Value   `json:"executor_id"`
	Command    Command `json:"command"`
}

// Command comment
type Command struct {
	Shell bool   `json:"shell"`
	Value string `json:"value"`
}

// HealthCheck comment
type HealthCheck struct {
	Type                string  `json:"type"`
	DelaySeconds        float64 `json:"delay_seconds"`
	IntervalSeconds     float64 `json:"interval_seconds"`
	TimeoutSeconds      float64 `json:"timeout_seconds"`
	ConsecutiveFailures int     `json:"consecutive_failures"`
	GracePeriodSeconds  float64 `json:"grace_period_seconds"`
	Command             Command `json:"command"`
}

// Status comment
type Status struct {
	UUID            string          `json:"uuid"`
	Timestamp       float64         `json:"timestamp"`
	AgentID         Value           `json:"agent_id"`
	ExecutorID      Value           `json:"executor_id"`
	TaskID          Value           `json:"task_id"`
	Source          string          `json:"source"`
	State           string          `json:"state"`
	ContainerStatus ContainerStatus `json:"container_status"`
}

// ContainerStatus comment
type ContainerStatus struct {
	ExecutorPID  int           `json:"executor_pid"`
	NetworkInfos []IPAddresses `json:"network_infos"`
}

// IPAddresses comment
type IPAddresses struct {
	IPAddresses []IPAddress `json:"ip_addresses"`
}

// IPAddress comment
type IPAddress struct {
	IPAddress string `json:"ip_address"`
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
