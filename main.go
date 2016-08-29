// GOlang implementation of mesos framework HTTP interface
//

package main

import (
	"log"
	"strconv"

	"github.com/tymofii-polekhin/megos-framework/scheduler"
	"github.com/tymofii-polekhin/megos-framework/utils"
)

// OfferCheck for resources
func OfferCheck(s *scheduler.Scheduler, o *utils.Offer) (accept bool) {
	return false
}

// BuildTaskInfo for new task
// func BuildTaskInfo(o *utils.Offer) (a *utils.Accept) {
//
// 	accept := new(utils.Accept)
//
// 	accept.OfferIDs = append(accept.OfferIDs, o.ID)
// 	accept.Filters.RefuseSeconds = 60.0
// 	accept.Operations = make([]utils.Operation, 1)
// 	accept.Operations[0].Type = "LAUNCH"
// 	accept.Operations[0].Launch.TaskInfos = make([]utils.TaskInfo, 1)
// 	accept.Operations[0].Launch.TaskInfos[0].Name = "megos-task-0"
// 	accept.Operations[0].Launch.TaskInfos[0].TaskID.Value = "12345-6789-12345-megos-task-0"
// 	accept.Operations[0].Launch.TaskInfos[0].AgentID = o.AgentID
// 	accept.Operations[0].Launch.TaskInfos[0].Command = new(utils.Command)
// 	accept.Operations[0].Launch.TaskInfos[0].Command.Shell = true
// 	accept.Operations[0].Launch.TaskInfos[0].Command.Value = "sleep 30"
//
// 	var cpus, mem, disk, ports utils.Resource
// 	cpus.Name = "cpus"
// 	cpus.Type = "SCALAR"
// 	cpus.Role = "*"
// 	cpus.Scalar = &utils.Scalar{Value: 0.1}
//
// 	mem.Name = "mem"
// 	mem.Type = "SCALAR"
// 	mem.Role = "*"
// 	mem.Scalar = &utils.Scalar{Value: 32.1}
//
// 	disk.Name = "disk"
// 	disk.Type = "SCALAR"
// 	disk.Role = "*"
// 	disk.Scalar = &utils.Scalar{Value: 128.1}
//
// 	ports.Name = "ports"
// 	ports.Type = "RANGES"
// 	ports.Role = "*"
// 	ports.Ranges = new([]utils.Ranges)
// 	ports.Ranges = make(utils.Ranges, 1)
// 	ports.Ranges.Range[0].Begin = 31000
// 	ports.Ranges.Range[0].End = 31000
//
// 	accept.Operations[0].Launch.TaskInfos[0].Resources = append(
// 		accept.Operations[0].Launch.TaskInfos[0].Resources, cpus, mem, disk)
//
// 	return accept
// }

func main() {

	// Desired workflow:
	//
	// s := new(scheduler.Scheduler)
	//
	// s.Subscribe(master, name, user)
	// defer s.Teardown()
	//
	// s.AddTask(task utils.TaskInfo)
	// s.AddTasks(tasks []utils.TaskInfo)
	//
	// s.TaskFitOffer(task, offer)
	// s.TasksFitOffer(tasks, offer)
	//
	// s.LaunchTask(task taskInfo, offer utils.Offer)
	// s.LaunchTasks(tasks []taskInfos, offers []utils.Offer, strategy)
	//
	// s.KillTask(task *TaskInfo)
	// s.KillTasks(tasks *[]utils.TaskInfo)
	//
	//

	// create new Scheduler
	s := new(scheduler.Scheduler)

	// subscribe
	err := s.Subscribe("localhost:5050", "MeGos Framework", "root")
	if err != nil {
		log.Fatalln(err)
	}
	// and unsibscribe later
	defer s.Teardown()

	// Creating task template
	task := new(utils.TaskInfo)
	task.AddCpus(0.1)
	task.AddMem(32.0)
	task.AddDisk(32.0)
	task.AddPort(31000)
	task.AddCommand("sleep 60")

	// Add 3 tasks with same resources
	for i := 0; i < 3; i++ {
		task.Name = "megos-task-" + strconv.Itoa(i)
		s.Tasks = append(s.Tasks, *task)
	}

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
			for _, offer := range event.Offers.Offers {
				offerUsed := false
				for _, task := range s.Tasks {
					if s.TaskIsLaunched(&task) == false {
						if task.FitOffer(&offer) {
							err := s.LaunchTask(&task, &offer)
							if err != nil {
								log.Println("Failed to launch task", task, "using offer", offer)
							} else {
								log.Println("Successfully submited task", task, "using offer", offer)
								break
							}
						} else { // task dont fit in offer
							log.Println("Offer", offer, "is not suitable for task", task)
						}
					} else { // task is already launched
						log.Println("Task is already launched")
					}
				} // for tasks
				for _, t := range s.LaunchedTasks {
					if t.AgentID == offer.AgentID {
						offerUsed = true
					}
				}
				if offerUsed != true {
					s.Decline(&offer, 5)
				}
			} // for offers

		case "UPDATE":
			log.Println("Received task", event.Update.Status.TaskID.Value,
				"status update", event.Update.Status.State)
			if event.Update.Status.State != "TASK_RUNNING" {
			}
			s.Acknowledge(event)

		case "HEARTBEAT":
			log.Println("Received heartbeat event")

		default:
			log.Fatalln("Unknown event type!")

		}

	} //forever loop
}
