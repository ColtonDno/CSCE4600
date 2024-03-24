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
		RemainingTime  int64
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
		i               int64
		currentProcess  *Process
		previousProcess *Process
		queue           = make([]*Process, 0)
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	for i = 0; len(queue) > 0 || i == 0; i++ {
		//Add processes to the queue at their arrival time
		if i < int64(len(processes)) {
			for j := range processes {
				if processes[j].ArrivalTime == int64(i) {
					queue = append(queue, &processes[j])
				}
			}
		}

		//Sort the queue by remaining run time
		slices.SortFunc(queue,
			func(a, b *Process) int {
				return cmp.Compare(a.RemainingTime, b.RemainingTime)
			})
		currentProcess = queue[0]

		//Update process times
		currentProcess.RemainingTime--
		for j := 1; j < len(queue); j++ {
			queue[j].WaitTime++
		}

		//When a process' run time reaches 0, set its completion time and remove it from the queue
		for j := len(queue) - 1; j >= 0; j-- {
			if queue[j].RemainingTime == 0 {
				queue[j].CompletionTime = i - 1
				queue = append(queue[:j], queue[j+1:]...)
			}
		}

		if i != 0 && currentProcess.ProcessID != previousProcess.ProcessID {
			gantt = append(gantt, TimeSlice{
				PID:   previousProcess.ProcessID,
				Start: previousProcess.StartTime,
				Stop:  i,
			})

			currentProcess.StartTime = i
		}

		previousProcess = currentProcess
	}

	//Add the final process to the list
	gantt = append(gantt, TimeSlice{
		PID:   currentProcess.ProcessID,
		Start: currentProcess.StartTime,
		Stop:  int64(i),
	})

	for i := range processes {
		var process *Process = &processes[i]
		process.TurnaroundTime = process.BurstDuration + process.WaitTime
		process.CompletionTime = process.TurnaroundTime + process.ArrivalTime

		totalWait += float64(process.WaitTime)
		totalTurnaround += float64(process.TurnaroundTime)
		lastCompletion += float64(process.CompletionTime)

		schedule[i] = []string{
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
		i               int64
		currentProcess  *Process
		previousProcess *Process
		queue           = make([]*Process, 0)
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	for i = 0; len(queue) > 0 || i == 0; i++ {
		//Add processes to the queue at their arrival time
		if i < int64(len(processes)) {
			for j := range processes {
				if processes[j].ArrivalTime == int64(i) {
					queue = append(queue, &processes[j])
				}
			}
		}

		//Sort the queue by remaining run time
		slices.SortFunc(queue,
			func(a, b *Process) int {
				return cmp.Compare(a.RemainingTime, b.RemainingTime)
			})

		//Sort the queue by priority
		slices.SortFunc(queue,
			func(a, b *Process) int {
				return cmp.Compare(a.Priority, b.Priority)
			})
		currentProcess = queue[0]

		//Update process times
		currentProcess.RemainingTime--
		for j := 1; j < len(queue); j++ {
			queue[j].WaitTime++
		}

		/*fmt.Println(i)
		if i > 0 {
			fmt.Print(previousProcess.ProcessID)
			fmt.Print(" -> ")
		}
		fmt.Println(currentProcess.ProcessID)
		fmt.Println()*/

		//When a process' run time reaches 0, set its completion time and remove it from the queue
		for j := len(queue) - 1; j >= 0; j-- {
			if queue[j].RemainingTime == 0 {
				queue[j].CompletionTime = i - 1
				queue = append(queue[:j], queue[j+1:]...)
			}
		}

		if i != 0 && currentProcess.ProcessID != previousProcess.ProcessID {
			gantt = append(gantt, TimeSlice{
				PID:   previousProcess.ProcessID,
				Start: previousProcess.StartTime,
				Stop:  i,
			})

			currentProcess.StartTime = i
		}

		previousProcess = currentProcess
	}

	//Add the final process to the list
	gantt = append(gantt, TimeSlice{
		PID:   currentProcess.ProcessID,
		Start: currentProcess.StartTime,
		Stop:  int64(i),
	})

	for i := range processes {
		var process *Process = &processes[i]
		process.TurnaroundTime = process.BurstDuration + process.WaitTime
		process.CompletionTime = process.TurnaroundTime + process.ArrivalTime

		totalWait += float64(process.WaitTime)
		totalTurnaround += float64(process.TurnaroundTime)
		lastCompletion += float64(process.CompletionTime)

		schedule[i] = []string{
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
		i               int64
		//timeQuantum     int64 = 4
		currentProcess  *Process
		previousProcess *Process
		queue           = make([]*Process, 0)
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	for i = 0; len(queue) > 0 || i == 0; i++ {
		//Add processes to the queue at their arrival time
		if i < int64(len(processes)) {
			for j := range processes {
				if processes[j].ArrivalTime == int64(i) {
					queue = append(queue, &processes[j])
				}
			}
		}
		currentProcess = queue[0]

		//Update process times
		currentProcess.RemainingTime--
		for j := 1; j < len(queue); j++ {
			queue[j].WaitTime++
		}

		//When a process' run time reaches 0, set its completion time and remove it from the queue
		for j := len(queue) - 1; j >= 0; j-- {
			if queue[j].RemainingTime == 0 {
				queue[j].CompletionTime = i - 1
				queue = append(queue[:j], queue[j+1:]...)
			}
		}

		if len(queue) > 1 && i-currentProcess.StartTime == 4 {
			fmt.Println(queue)
			queue = append(queue[1:], queue[0])
			fmt.Println(queue)
		}

		if i != 0 && currentProcess.ProcessID != previousProcess.ProcessID {
			gantt = append(gantt, TimeSlice{
				PID:   previousProcess.ProcessID,
				Start: previousProcess.StartTime,
				Stop:  i,
			})

			currentProcess.StartTime = i
		}

		previousProcess = currentProcess
	}

	//Add the final process to the list
	gantt = append(gantt, TimeSlice{
		PID:   currentProcess.ProcessID,
		Start: currentProcess.StartTime,
		Stop:  int64(i),
	})

	for i := range processes {
		var process *Process = &processes[i]
		process.TurnaroundTime = process.BurstDuration + process.WaitTime
		process.CompletionTime = process.TurnaroundTime + process.ArrivalTime

		totalWait += float64(process.WaitTime)
		totalTurnaround += float64(process.TurnaroundTime)
		lastCompletion += float64(process.CompletionTime)

		schedule[i] = []string{
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
