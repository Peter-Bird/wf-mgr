package handlers

import (
	"net/http"

	"fmt"
	"wf-mgr/services"

	mdl "github.com/Peter-Bird/models"
)

// ExecuteWorkflowHandler handles workflow execution
func ExecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	workflow := generateWorkFlow()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Send initial response
	fmt.Fprintf(w, "Initial response\n")
	flusher.Flush()

	fmt.Fprintf(w, "Workflow execution started\n")
	flusher.Flush()

	services.ExecuteWorkflow(w, workflow)
}

func generateWorkFlow() mdl.Workflow {

	steps := []mdl.Step{
		{
			Endpoint:     "http://localhost:8081/workflows",
			Method:       "GET",
			Parameters:   nil,
			Dependencies: []string{},
		},
	}

	parameters := map[string]interface{}{
		"id":    "WF1",
		"name":  "WF1",
		"steps": steps,
	}

	return mdl.Workflow{
		Id:   "WORKFLOW_ID",
		Name: "Sample Workflow",
		Steps: []mdl.Step{
			{
				Endpoint:     "http://localhost:8081/workflows",
				Method:       "GET",
				Parameters:   nil,
				Dependencies: nil,
			},
			{
				Endpoint:     "http://localhost:8081/workflows/submit",
				Method:       "POST",
				Parameters:   parameters,
				Dependencies: []string{"step1"},
			},
			{
				Endpoint:     "http://localhost:8081/workflows/get/WF1",
				Method:       "GET",
				Parameters:   nil,
				Dependencies: []string{"step2"},
			},
		},
	}
}
