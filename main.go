// GOlang implementation of mesos framework HTTP interface
//

package main

import (
	"log"

	"github.com/tymofii-polekhin/megos-framework/scheduler"
	"github.com/tymofii-polekhin/megos-framework/utils"
)

// ProcessOffer comment
func ProcessOffer(s *scheduler.Scheduler, o *utils.Offer) (err error) {

	if s.TasksLaunched < s.TasksTotal { //accept offer

		var accept utils.Accept
		accept.Type = "ACCEPT"
		accept.FrameworkID.Value = s.FrameworkID

		accept.Accept.OfferIDs = append(accept.Accept.OfferIDs, o.ID)

		var taskinfo utils.TaskInfo
		taskinfo.Name = "Go Sample Task"
		taskinfo.TaskID = utils.Value{Value: "12345-6789-54321-go-sample-task"}
		taskinfo.AgentID = o.AgentID
		taskinfo.Executor.ExecutorID = utils.Value{Value: "12345-6789-54321-go-sample-executor"}
		taskinfo.Executor.Command.Shell = true
		taskinfo.Executor.Command.Value = "sleep 600"

		var taskcpu, taskmem utils.Resource

		taskcpu.Name = "cpus"
		taskcpu.Role = "*"
		taskcpu.Type = "SCALAR"
		taskcpu.Scalar = &utils.Scalar{Value: 0.5}

		taskmem.Name = "mem"
		taskmem.Role = "*"
		taskmem.Type = "SCALAR"
		taskmem.Scalar = &utils.Scalar{Value: 128.0}

		taskinfo.Resources = append(taskinfo.Resources, taskcpu, taskmem)

		var task utils.Operation
		task.Type = "LAUNCH"
		task.Launch.TaskInfos = append(task.Launch.TaskInfos, taskinfo)

		accept.Accept.Operations = append(accept.Accept.Operations, task)

		accept.Accept.Filters.RefuseSeconds = 5.0

		// send accepted offer to channel
		log.Println("Accepted offer:", o)
		//s.AcceptedOffers <- accept
		r, err := s.AcceptOffer(accept)
		if err != nil {
			log.Fatalln(err)
		}
		if (err == nil) && (r != "") {
			s.TasksLaunched++
			log.Println("Response from master:", r)
		}

	} else { //decline offer

		log.Println("TasksLaunched:", s.TasksLaunched, "TasksTotal:", s.TasksTotal)
		log.Println("Declining offer:", o)
		aknowleged, err := s.DeclineOffer(o.ID.Value)
		if err != nil {
			log.Fatal(err)
		}
		if (err == nil) && (aknowleged != "") {
			log.Println("Response from master:", aknowleged)
		}

	}

	return nil
}

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

	// we want to launch one task
	s.TasksTotal = 1

	for {
		select {
		case offer := <-s.NewOffers:
			log.Println("Offer received:", offer)
			ProcessOffer(s, &offer)
		}
	} //forever loop
}
