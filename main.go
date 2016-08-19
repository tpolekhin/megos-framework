// GOlang implementation of mesos framework HTTP interface
//

package main

import (
	"log"

	"github.com/tymofii-polekhin/megos-framework/scheduler"
)

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
