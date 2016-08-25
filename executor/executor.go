package executor

import "github.com/tymofii-polekhin/megos-framework/utils"

// Executor info
type Executor struct {
	ExecutorName  string
	ExecutorID    utils.Value
	AgentID       utils.Value
	AgentEndpoint string
}

// Subscribe executor to corresponding Agent endpoint
func (e *Executor) Subscribe() (err error) {

	return nil
}
