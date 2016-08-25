package utils

// Events: Master -> Scheduler
// SUBSCRIBED
// OFFERS
// RESCIND
// UPDATE
// MESSAGE
// FAILURE
// ERROR
// HEARTBEAT

// Event message from mesos to scheduler/executor
type Event struct {
	Type       string      `json:"type"`
	Subscribed *subscribed `json:"subscribed,omitempty"`
	Offers     *offers     `json:"offers,omitempty"`
	Recind     *recind     `json:"recind,omitempty"`
	Update     *update     `json:"update,omitempty"`
	Message    *message    `json:"message,omitempty"`
	Failure    *failure    `json:"failure,omitempty"`
	Error      *merror     `json:"error,omitempty"`
}

type subscribed struct {
	HeartbeatIntervalSeconds float64 `json:"heartbeat_interval_seconds"`
	FrameworkID              Value   `json:"framework_id"`
}

type offers struct {
	Offers []Offer `json:"offers"`
}

type recind struct{}

type update struct {
	Status Status `json:"status"`
}

type message struct{}
type failure struct{}
type merror struct{}
