package main

import (
	"fmt"
	"os"

	"github.com/mackerelio/go-osstat/disk"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/shirou/gopsutil/cpu"
	v8 "rogchap.com/v8go"
)

type CpuDetail struct {
	ModelName string
	CPU       int32
}

func main() {

	runJs()
	//__cpu()
	// _disk()
	// _cpu()
	// mem()
}

func runJs() {
	ctx := v8.NewContext()                                  // creates a new V8 context with a new Isolate aka VM
	ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
	ctx.RunScript("const result = add(3, 4)", "main.js")    // any functions previously added to the context can be called
	val, _ := ctx.RunScript("result", "value.js")           // return a value in JavaScript back to Go
	fmt.Printf("addition result: %s", val)
}

func mem() {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	fmt.Println("memory total: ", ByteCountSI(memory.Total))
	fmt.Println("memory used: ", ByteCountSI(memory.Used))
	fmt.Println("memory cached: ", ByteCountSI(memory.Cached))
	fmt.Println("memory free: ", ByteCountSI(memory.Free))
}

func __cpu() {
	cpuStat, err := cpu.Info()
	if err != nil {
		fmt.Println("error", err)
		return
	}

	output := make([]CpuDetail, len(cpuStat))

	percentage, errP := cpu.Percent(0, true)
	if errP != nil {
		fmt.Println("error", errP)
		return
	}

	fmt.Println("CPU Percent ", percentage)

	var i int = 0
	for _, val := range cpuStat {
		c := CpuDetail{ModelName: val.ModelName, CPU: val.CPU}
		output[i] = c
		fmt.Println("CPU Info \n", val)
	}
}

// func _cpu() {
// 	before, err := cpu.Get()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return
// 	}
// 	time.Sleep(time.Duration(1) * time.Second)
// 	after, err := cpu.Get()
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		return
// 	}
// 	total := float64(after.Total - before.Total)
// 	fmt.Printf("cpu used: %f %%\n", float64(after.User-before.User)/total*100)
// 	fmt.Printf("cpu system: %f %%\n", float64(after.System-before.System)/total*100)
// 	fmt.Printf("cpu idle: %f %%\n", float64(after.Idle-before.Idle)/total*100)
// }

func _disk() {
	disk, err := disk.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	fmt.Printf("disk name: %v", disk)

	for _, val := range disk {
		fmt.Printf("disk name: %v reads completed: %v write completed: %v \n", val.Name, val.ReadsCompleted, val.WritesCompleted)
	}
}

func ByteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
