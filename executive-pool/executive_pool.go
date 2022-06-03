package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
)

const tempDir = "templates"

type ExecutivePoolServer struct {
	port    uint16
	gateway string
}

type RunnableFunction struct {
	JsFunction *string `json:"runnable_function"`
}

type ResultResponse struct {
	Result []string `json:"result"`
}

func (ep *ExecutivePoolServer) Port() uint16 {
	return ep.port
}

func (ep *ExecutivePoolServer) Gateway() string {
	return ep.gateway
}

func NewExecutivePoolServer(port uint16, gateway string) *ExecutivePoolServer {
	return &ExecutivePoolServer{port, gateway}
}

func JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}

func (r *RunnableFunction) Validate() bool {
	if r.JsFunction == nil {
		return false
	}
	return true
}

func (ep *ExecutivePoolServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
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

		bt := RunnableFunction{runnableFunction.JsFunction}
		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)
		resp, _ := http.Post(ep.Gateway()+"/compute-request", "application/json", buf)
		if resp.StatusCode == 200 {
			decoder := json.NewDecoder(resp.Body)
			var bar ResultResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(JsonStatus("fail")))
				return
			}

			m, _ := json.Marshal(struct {
				ComputeType string   `json:"compute_type"`
				Result      []string `json:"process_result"`
				Status      string   `json:"status"`
			}{
				ComputeType: "V8 JavaScript engine",
				Result:      bar.Result,
				Status:      "success",
			})
			io.WriteString(w, string(m[:]))
			return
		}
		io.WriteString(w, string(JsonStatus("fail")))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (ep *ExecutivePoolServer) Run() {
	http.HandleFunc("/", ep.Index)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ep.Port())), nil))
}
