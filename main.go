// GOlang implementation of mesos framework HTTP interface
//

package main

import (
	"log"

	"github.com/tymofii-polekhin/megos-framework/scheduler"
	"github.com/tymofii-polekhin/megos-framework/utils"
)

// OfferCheck for resources
func OfferCheck(s *scheduler.Scheduler, o *utils.Offer) (accept bool) {
	if s.TasksLaunched < s.TasksTotal {
		return true
	}
	return false
}

// BuildTaskInfo for new task
func BuildTaskInfo(o *utils.Offer) (a *utils.Accept) {

	accept := new(utils.Accept)

	accept.OfferIDs = append(accept.OfferIDs, o.ID)
	accept.Filters.RefuseSeconds = 60.0
	accept.Operations = make([]utils.Operation, 1)
	accept.Operations[0].Type = "LAUNCH"
	accept.Operations[0].Launch.TaskInfos = make([]utils.TaskInfo, 1)
	accept.Operations[0].Launch.TaskInfos[0].Name = "megos-task-0"
	accept.Operations[0].Launch.TaskInfos[0].TaskID.Value = "12345-6789-12345-megos-task-0"
	accept.Operations[0].Launch.TaskInfos[0].AgentID = o.AgentID
	accept.Operations[0].Launch.TaskInfos[0].Command = new(utils.Command)
	accept.Operations[0].Launch.TaskInfos[0].Command.Shell = true
	accept.Operations[0].Launch.TaskInfos[0].Command.Value = "sleep 30"

	var cpus, mem, disk utils.Resource
	cpus.Name = "cpus"
	cpus.Type = "SCALAR"
	cpus.Role = "*"
	cpus.Scalar = &utils.Scalar{Value: 0.1}

	mem.Name = "mem"
	mem.Type = "SCALAR"
	mem.Role = "*"
	mem.Scalar = &utils.Scalar{Value: 32.1}

	disk.Name = "disk"
	disk.Type = "SCALAR"
	disk.Role = "*"
	disk.Scalar = &utils.Scalar{Value: 128.1}

	// ports.Name = "ports"
	// ports.Type = "RANGES"
	// ports.Role = "*"
	// ports.Ranges = new([]utils.Ranges)
	// ports.Ranges = make(utils.Ranges, 1)
	// ports.Ranges.Range[0].Begin = 31000
	// ports.Ranges.Range[0].End = 31000

	accept.Operations[0].Launch.TaskInfos[0].Resources = append(
		accept.Operations[0].Launch.TaskInfos[0].Resources, cpus, mem, disk)

	return accept
}

func main() {

	// create new Scheduler
	s := new(scheduler.Scheduler)

	// subscribe
	err := s.Subscribe("localhost:5050", "MeGos Framework", "root")
	if err != nil {
		log.Fatalln(err)
	}
	// and unsibscribe later
	defer s.Teardown()

	// we want to launch one task
	s.TasksLaunched = 0
	s.TasksTotal = 1

	for {
		event := <-s.Events

		switch event.Type {

		case "SUBSCRIBED":
			s.FrameworkID = event.Subscribed.FrameworkID
			s.HeartbeatIntervalSeconds = event.Subscribed.HeartbeatIntervalSeconds
			log.Println("Subscribed to", s.Master, "with frameworkID", s.FrameworkID,
				"and heartbeat interval of", s.HeartbeatIntervalSeconds, "seconds")

		case "OFFERS":
			log.Println("Received", len(event.Offers.Offers), "offers from master")
			var decline []utils.Value
			for _, o := range event.Offers.Offers {
				if OfferCheck(s, &o) {
					err := s.Accept(BuildTaskInfo(&o))
					if err != nil {
						log.Println(err)
						decline = append(decline, o.ID)
					} else {
						log.Println("Launched task of offer", o.ID.Value)
						s.TasksLaunched = 1
					}
				} else {
					decline = append(decline, o.ID)
				}
			}
			s.Decline(decline, 5)

		case "UPDATE":
			log.Println("Received task", event.Update.Status.TaskID.Value,
				"status update", event.Update.Status.State)
			if event.Update.Status.State != "TASK_RUNNING" {
				s.TasksLaunched = 0
			}
			s.Acknowledge(event)

		case "HEARTBEAT":
			log.Println("Received heartbeat event")

		default:
			log.Fatalln("Unknown event type!")

		}

	} //forever loop
}
