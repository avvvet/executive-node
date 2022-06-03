package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	v8 "rogchap.com/v8go"
)

type ExecutiveNode struct {
	port uint16
}

type RunnableFunction struct {
	JsFunction *string `json:"runnable_function"`
}

type CpuInfo struct {
	ModelName string
	VendorID  string
	Model     string
	Mhz       float64
	CPU       int32
}

type ResultResponse struct {
	Result []string `json:"result"`
}

var ComputeLog []string

func (ar *ResultResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Result []string `json:"result"`
	}{
		Result: ar.Result,
	})
}

func (c *CpuInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ModelName string  `json:"model_name"`
		VendorID  string  `json:"vendor_id"`
		Model     string  `json:"model"`
		Mhz       float64 `json:"mhz"`
		CPU       int32   `json:"cpu"`
	}{
		ModelName: c.ModelName,
		VendorID:  c.VendorID,
		Model:     c.Model,
		Mhz:       c.Mhz,
		CPU:       c.CPU,
	})
}

func NewExecutiveNode(port uint16) *ExecutiveNode {
	return &ExecutiveNode{port}
}

func (s *ExecutiveNode) Port() uint16 {
	return s.port
}

func (r *RunnableFunction) Validate() bool {
	if r.JsFunction == nil {
		return false
	}
	return true
}

func (s *ExecutiveNode) CpuInfo(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		cpuInfo := cpuInfo()
		m, _ := json.Marshal(struct {
			CpuInfos []CpuInfo `json:"cpu_info"`
			Length   int       `json:"length"`
		}{
			CpuInfos: cpuInfo,
			Length:   len(cpuInfo),
		})
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *ExecutiveNode) CpuPercent(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		cpuPercent := cpuPercent()
		m, _ := json.Marshal(struct {
			CpuPercent []float64 `json:"cpu_info"`
			Length     int       `json:"length"`
		}{
			CpuPercent: cpuPercent,
			Length:     len(cpuPercent),
		})
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *ExecutiveNode) ComputeRequest(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		n := rand.Intn(500000)
		result := factorial(n)
		m, _ := json.Marshal(struct {
			ComputeType string `json:"compute_type"`
			Result      int    `json:"process_result"`
			Status      string `json:"status"`
		}{
			ComputeType: "Graph travers",
			Result:      result,
			Status:      "success",
		})
		io.WriteString(w, string(m[:]))
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var runnableFunction RunnableFunction
		err := decoder.Decode(&runnableFunction)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(JsonStatus("runnable fail")))
			return
		}
		if !runnableFunction.Validate() {
			log.Println("ERROR: missing runnable function field(s)")
			io.WriteString(w, string(JsonStatus("runnable function fail")))
			return
		}
		fn := *runnableFunction.JsFunction

		w.Header().Add("Content-Type", "application/json")

		result := V8engine(fn)

		ar := &ResultResponse{result}
		m, _ := ar.MarshalJSON()

		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP Method  *****")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *ExecutiveNode) ProcessLogs(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")

		m, _ := json.Marshal(struct {
			ProcessLogs []string `json:"compute_logs"`
			Length      int      `json:"length"`
		}{
			ProcessLogs: ComputeLog,
			Length:      len(ComputeLog),
		})
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *ExecutiveNode) Run() {
	http.HandleFunc("/", s.CpuInfo)
	http.HandleFunc("/state", s.CpuPercent)
	http.HandleFunc("/compute-request", s.ComputeRequest)
	http.HandleFunc("/logs", s.ProcessLogs)
	log.Println("Executive node running at localhost:" + strconv.Itoa(int(s.Port())) + " ...")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))
}

func JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}

func cpuInfo() []CpuInfo {
	cpuStat, err := cpu.Info()
	if err != nil {
		fmt.Println("error", err)
		return nil
	}

	output := make([]CpuInfo, len(cpuStat))
	var i int = 0
	for _, v := range cpuStat {
		c := CpuInfo{ModelName: v.ModelName, VendorID: v.VendorID, Model: v.Model, Mhz: v.Mhz, CPU: v.CPU}
		output[i] = c
		i += 1
	}

	return output
}

func cpuPercent() []float64 {
	percentage, err := cpu.Percent(0, true)
	if err != nil {
		fmt.Println("error", err)
		return nil
	}
	return percentage
}

func factorial(num int) int {
	defer timeTrack(time.Now(), "graph traversal")
	//in seconds to sleep.
	var r = rand.Intn(10)
	time.Sleep(time.Duration(r) * time.Second)

	return num
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	ComputeLog = append(ComputeLog, "compute took: "+elapsed.String()+" compute type: "+name)

	log.Printf("%s compute took: %s", name, elapsed)
}

func V8engine(runnableFn string) []string {
	defer timeTrack(time.Now(), "v8 engine JavaScript")
	ctx := v8.NewContext()               // creates a new V8 context with a new Isolate aka VM
	ctx.RunScript(runnableFn, "math.js") // executes a script on the global context
	//ctx.RunScript("const result = add(3, 4)", "main.js") // any functions previously added to the context can be called
	val, _ := ctx.RunScript("result", "value.js") // return a value in JavaScript back to Go
	fmt.Printf("addition result: %s", val)
	var r []string
	r = append(r, val.String())
	return r
}
