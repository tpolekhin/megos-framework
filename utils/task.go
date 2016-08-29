package utils

import "fmt"

const scalar = "SCALAR"
const ranges = "RANGES"

// TaskInfo comment
type TaskInfo struct {
	Name        string       `json:"name"`
	TaskID      Value        `json:"task_id"`
	AgentID     Value        `json:"agent_id"`
	Command     *Command     `json:"command,omitempty"`
	HealthCheck *HealthCheck `json:"health_check,omitempty"`
	Executor    *Executor    `json:"executor,omitempty"`
	Resources   []Resource   `json:"resources"`
}

// AddCommand to execute
func (task *TaskInfo) AddCommand(command string) {
	task.Command = new(Command)
	task.Command.Shell = true
	task.Command.Value = command
}

// AddCpus to task resources
func (task *TaskInfo) AddCpus(cpus float64) (err error) {
	if cpus > 0.0 {
		r := new(Resource)
		r.Name = "cpus"
		r.Role = "*"
		r.Type = scalar
		r.Scalar = &Scalar{Value: cpus}
		task.Resources = append(task.Resources, *r)
	} else {
		return fmt.Errorf("CPUs value cannot be less than zero")
	}
	return nil
}

// AddMem to task resources
func (task *TaskInfo) AddMem(mem float64) (err error) {
	if mem > 0.0 {
		r := new(Resource)
		r.Name = "mem"
		r.Role = "*"
		r.Type = scalar
		r.Scalar = &Scalar{Value: mem}
		task.Resources = append(task.Resources, *r)
	} else {
		return fmt.Errorf("MEM value cannot be less than zero")
	}
	return nil
}

// AddDisk to task resources
func (task *TaskInfo) AddDisk(disk float64) (err error) {
	if disk > 0.0 {
		r := new(Resource)
		r.Name = "disk"
		r.Role = "*"
		r.Type = scalar
		r.Scalar = &Scalar{Value: disk}
		task.Resources = append(task.Resources, *r)
	} else {
		return fmt.Errorf("Disk value cannot be less than zero")
	}
	return nil
}

// AddPort to task resources
func (task *TaskInfo) AddPort(port int) (err error) {
	if port > 0 {
		r := new(Resource)
		r.Name = "ports"
		r.Role = "*"
		r.Type = ranges
		r.Ranges = new(Ranges)
		r.Ranges.Range = make([]Range, 1)
		r.Ranges.Range[0].Begin = port
		r.Ranges.Range[0].End = port

		task.Resources = append(task.Resources, *r)
	} else {
		return fmt.Errorf("Port value cannot be less than zero")
	}
	return nil
}

func compareRanges(t, o *Ranges) (fit bool) {

	var fitting []bool // array of ranges fit

	if len(t.Range) > 0 { // we have at least one port range to fit
		if len(o.Range) > 0 { // we have at least one port range in offer

			fitting = make([]bool, len(t.Range)) // equal to task ranges

			// iterating task ranges trying to fit one of offer ranges
			for taskPortRangeIndex, taskPortRange := range t.Range {
				for _, offerPortRange := range o.Range {

					// if task range fit offer range
					if (taskPortRange.Begin >= offerPortRange.Begin) &&
						(taskPortRange.End <= offerPortRange.End) {
						fitting[taskPortRangeIndex] = true
					} // if task range fit offer range

				} // for offer ranges
			} // for task ranges

		} else {
			return false // we have no ports available in offer at all
		}
	} else {
		return true // we're not using ports, fit any offer
	}

	for _, fit := range fitting { // we must fit all task ranges
		if !fit {
			return false
		}
	}

	return true // we fit all task ranges
}

// FitOffer check offer resources to fit a task
func (task *TaskInfo) FitOffer(offer *Offer) (fit bool) {
	var fitcpus, fitmem, fitdisk, fitports bool

	for _, or := range offer.Resources {

		switch or.Name {

		case "cpus":
			for _, tr := range task.Resources {
				if tr.Name == "cpus" {
					if tr.Scalar.Value < or.Scalar.Value {
						fitcpus = true
					} else {
						fitcpus = false
					}
				}
			}

		case "mem":
			for _, tr := range task.Resources {
				if tr.Name == "mem" {
					if tr.Scalar.Value < or.Scalar.Value {
						fitmem = true
					} else {
						fitmem = false
					}
				}
			}

		case "disk":
			for _, tr := range task.Resources {
				if tr.Name == "disk" {
					if tr.Scalar.Value < or.Scalar.Value {
						fitdisk = true
					} else {
						fitdisk = false
					}
				}
			}

		case "ports":
			for _, tr := range task.Resources {
				if tr.Name == "ports" {
					if compareRanges(tr.Ranges, or.Ranges) {
						fitports = true
					} else {
						fitports = false
					}
				}
			}

		} // switch
	} // for

	if (fitcpus == true) && (fitmem == true) && (fitdisk == true) && (fitports == true) {
		return true
	}

	return false
}
