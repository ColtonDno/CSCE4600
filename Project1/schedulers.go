package main

import (
	"cmp"
	"fmt"
	"io"
	"slices"
)

type (
	Process struct {
		ProcessID      string
		ArrivalTime    int64
		BurstDuration  int64
		Priority       int64
		TimeRemaining  int64
		StartTime      int64
		WaitTime       int64
		TurnaroundTime int64
		CompletionTime int64
	}
	TimeSlice struct {
		PID   string
		Start int64
		Stop  int64
	}
)

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFSchedule(w io.Writer, title string, processes []Process) {
	var (
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		time_step       int64
		currentProcess  *Process
		previousProcess *Process
		queue           = make([]*Process, 0)
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	//Increase the time until the process queue is empty
	for time_step = 0; len(queue) > 0 || time_step == 0; time_step++ {
		//Add processes to the queue at their arrival time
		for j := range processes {
			if processes[j].ArrivalTime == int64(time_step) {
				queue = append(queue, &processes[j])

				//Sort the queue by remaining run time
				slices.SortFunc(queue,
					func(a, b *Process) int {
						return cmp.Compare(a.TimeRemaining, b.TimeRemaining)
					})
			}
		}

		currentProcess = queue[0]

		//Update process times
		currentProcess.TimeRemaining--
		for j := 1; j < len(queue); j++ {
			queue[j].WaitTime++
		}

		//When a process' run time reaches 0, set its completion time and remove it from the queue
		for j := len(queue) - 1; j >= 0; j-- {
			if queue[j].TimeRemaining == 0 {
				queue[j].CompletionTime = time_step - 1
				queue = append(queue[:j], queue[j+1:]...)
			}
		}

		//Add process to the gantt chart whenever a context switch occurs
		if time_step != 0 && currentProcess.ProcessID != previousProcess.ProcessID {
			gantt = append(gantt, TimeSlice{
				PID:   previousProcess.ProcessID,
				Start: previousProcess.StartTime,
				Stop:  time_step,
			})

			currentProcess.StartTime = time_step
		}

		previousProcess = currentProcess
	}

	//Add the final process to the list
	gantt = append(gantt, TimeSlice{
		PID:   currentProcess.ProcessID,
		Start: currentProcess.StartTime,
		Stop:  int64(time_step),
	})

	//Add processes to the schedule table
	for time_step := range processes {
		var process *Process = &processes[time_step]
		process.TurnaroundTime = process.BurstDuration + process.WaitTime
		process.CompletionTime = process.TurnaroundTime + process.ArrivalTime

		totalWait += float64(process.WaitTime)
		totalTurnaround += float64(process.TurnaroundTime)
		lastCompletion += float64(process.CompletionTime)

		schedule[time_step] = []string{
			fmt.Sprint(process.ProcessID),
			fmt.Sprint(process.Priority),
			fmt.Sprint(process.BurstDuration),
			fmt.Sprint(process.ArrivalTime),
			fmt.Sprint(process.WaitTime),
			fmt.Sprint(process.TurnaroundTime),
			fmt.Sprint(process.CompletionTime),
		}
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
	var (
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		time_step       int64
		currentProcess  *Process
		previousProcess *Process
		queue           = make([]*Process, 0)
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	//Increase the time until the process queue is empty
	for time_step = 0; len(queue) > 0 || time_step == 0; time_step++ {
		//Add processes to the queue at their arrival time
		for j := range processes {
			if processes[j].ArrivalTime == int64(time_step) {
				queue = append(queue, &processes[j])

				//Sort the queue by remaining run time
				slices.SortFunc(queue,
					func(a, b *Process) int {
						return cmp.Compare(a.TimeRemaining, b.TimeRemaining)
					})

				//Sort the queue by priority
				slices.SortFunc(queue,
					func(a, b *Process) int {
						return cmp.Compare(a.Priority, b.Priority)
					})
			}
		}

		currentProcess = queue[0]

		//Update process times
		currentProcess.TimeRemaining--
		for j := 1; j < len(queue); j++ {
			queue[j].WaitTime++
		}

		//When a process' run time reaches 0, set its completion time and remove it from the queue
		for j := len(queue) - 1; j >= 0; j-- {
			if queue[j].TimeRemaining == 0 {
				queue[j].CompletionTime = time_step - 1
				queue = append(queue[:j], queue[j+1:]...)
			}
		}

		//Add process to the gantt chart whenever a context switch occurs
		if time_step != 0 && currentProcess.ProcessID != previousProcess.ProcessID {
			gantt = append(gantt, TimeSlice{
				PID:   previousProcess.ProcessID,
				Start: previousProcess.StartTime,
				Stop:  time_step,
			})

			currentProcess.StartTime = time_step
		}

		previousProcess = currentProcess
	}

	//Add the final process to the list
	gantt = append(gantt, TimeSlice{
		PID:   currentProcess.ProcessID,
		Start: currentProcess.StartTime,
		Stop:  int64(time_step),
	})

	//Add processes to the schedule table
	for time_step := range processes {
		var process *Process = &processes[time_step]
		process.TurnaroundTime = process.BurstDuration + process.WaitTime
		process.CompletionTime = process.TurnaroundTime + process.ArrivalTime

		totalWait += float64(process.WaitTime)
		totalTurnaround += float64(process.TurnaroundTime)
		lastCompletion += float64(process.CompletionTime)

		schedule[time_step] = []string{
			fmt.Sprint(process.ProcessID),
			fmt.Sprint(process.Priority),
			fmt.Sprint(process.BurstDuration),
			fmt.Sprint(process.ArrivalTime),
			fmt.Sprint(process.WaitTime),
			fmt.Sprint(process.TurnaroundTime),
			fmt.Sprint(process.CompletionTime),
		}
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func RRSchedule(w io.Writer, title string, processes []Process) {
	var (
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		time_step       int64
		timeQuantum     int64 = 4
		currentProcess  *Process
		previousProcess *Process
		queue           = make([]*Process, 0)
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	//Increase the time until the process queue is empty
	for time_step = 0; len(queue) > 0 || time_step == 0; time_step++ {
		//Add processes to the queue at their arrival time
		for j := range processes {
			if processes[j].ArrivalTime == int64(time_step) {
				queue = append(queue, &processes[j])
			}
		}
		currentProcess = queue[0]

		//Update process times
		currentProcess.TimeRemaining--
		for j := 1; j < len(queue); j++ {
			queue[j].WaitTime++
		}

		//When a process' run time reaches 0, set its completion time and remove it from the queue
		for j := len(queue) - 1; j >= 0; j-- {
			if queue[j].TimeRemaining == 0 {
				queue[j].CompletionTime = time_step - 1
				queue = append(queue[:j], queue[j+1:]...)
			}
		}

		//Move process to the back of the queue if it exceeds the time quantum
		if len(queue) > 1 && time_step-currentProcess.StartTime == timeQuantum {
			fmt.Println(queue)
			queue = append(queue[1:], queue[0])
			fmt.Println(queue)
		}

		//Add process to the gantt chart whenever a context switch occurs
		if time_step != 0 && currentProcess.ProcessID != previousProcess.ProcessID {
			gantt = append(gantt, TimeSlice{
				PID:   previousProcess.ProcessID,
				Start: previousProcess.StartTime,
				Stop:  time_step,
			})

			currentProcess.StartTime = time_step
		}

		previousProcess = currentProcess
	}

	//Add the final process to the list
	gantt = append(gantt, TimeSlice{
		PID:   currentProcess.ProcessID,
		Start: currentProcess.StartTime,
		Stop:  int64(time_step),
	})

	//Add processes to the schedule table
	for time_step := range processes {
		var process *Process = &processes[time_step]
		process.TurnaroundTime = process.BurstDuration + process.WaitTime
		process.CompletionTime = process.TurnaroundTime + process.ArrivalTime

		totalWait += float64(process.WaitTime)
		totalTurnaround += float64(process.TurnaroundTime)
		lastCompletion += float64(process.CompletionTime)

		schedule[time_step] = []string{
			fmt.Sprint(process.ProcessID),
			fmt.Sprint(process.Priority),
			fmt.Sprint(process.BurstDuration),
			fmt.Sprint(process.ArrivalTime),
			fmt.Sprint(process.WaitTime),
			fmt.Sprint(process.TurnaroundTime),
			fmt.Sprint(process.CompletionTime),
		}
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

//endregion
